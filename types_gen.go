// +build ignore

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
)

var lengthFunction = `// Len return the length of this skipmap.
// Keep in sync with types_gen.go:lengthFunction
// Special case for code generation, Must in the tail of skipmap.go.
func (s *Int64Map) Len() int {
	return int(atomic.LoadInt64(&s.length))
}`

func main() {
	f, err := os.Open("skipmap.go")
	if err != nil {
		panic(err)
	}
	filedata, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	w := new(bytes.Buffer)
	w.WriteString(`// Code generated by go run types_gen.go; DO NOT EDIT.` + "\r\n")
	w.WriteString(string(filedata)[strings.Index(string(filedata), "package skipmap") : strings.Index(string(filedata), ")\n")+1])
	ts := []string{"Float32", "Float64", "Int32", "Int16", "Int", "Uint64", "Uint32", "Uint16", "Uint"} // all types need to be converted
	for _, upper := range ts {
		lower := strings.ToLower(upper)
		data := string(filedata)
		// Remove header.
		data = data[strings.Index(data, ")\n")+1:]
		// Remove the special case.
		data = strings.Replace(data, lengthFunction, "", -1)
		// Common cases.
		data = strings.Replace(data, "int64", lower, -1)
		data = strings.Replace(data, "Int64", upper, -1)
		if inSlice(lowerSlice(ts), lower) {
			data = strings.Replace(data, "length "+lower, "length int64", 1)
			data = strings.Replace(data, "atomic.Add"+upper, "atomic.AddInt64", -1)
		}
		// Add the special case.
		data = data + strings.Replace(lengthFunction, "Int64Map", upper+"Map", 1)
		w.WriteString(data)
		w.WriteString("\r\n")
	}

	// For desdending order.
	for _, upper := range append(ts, "Int64") {
		lower := strings.ToLower(upper)
		data := string(filedata)
		// Remove header.
		data = data[strings.Index(data, ")\n")+1:]
		// Remove the special case.
		data = strings.Replace(data, lengthFunction, "", -1)
		// DESC. (DIFF)
		data = strings.Replace(data, "ascending", "desdending", -1)
		data = strings.Replace(data, "NewInt64", "NewInt64Desc", -1)
		data = strings.Replace(data, "unlockInt64", "unlockInt64Desc", -1)
		data = strings.Replace(data, "Int64Map", "Int64MapDesc", -1)
		data = strings.Replace(data, "Int64Node", "Int64NodeDesc", -1)
		data = strings.Replace(data, "int64Node", "int64NodeDesc", -1)
		data = strings.Replace(data, "return n.key < key", "return n.key > key", -1)
		// Common cases.
		data = strings.Replace(data, "int64", lower, -1)
		data = strings.Replace(data, "Int64", upper, -1)
		if inSlice(lowerSlice(ts), lower) {
			data = strings.Replace(data, "length "+lower, "length int64", 1)
			data = strings.Replace(data, "atomic.Add"+upper, "atomic.AddInt64", -1)
		}
		// Add the special case. (DIFF)
		data = data + strings.Replace(lengthFunction, "Int64Map", upper+"MapDesc", 1)
		w.WriteString(data)
		w.WriteString("\r\n")
	}

	out, err := format.Source(w.Bytes())
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("types.go", out, 0660); err != nil {
		panic(err)
	}
}

func lowerSlice(s []string) []string {
	n := make([]string, len(s))
	for i, v := range s {
		n[i] = strings.ToLower(v)
	}
	return n
}

func inSlice(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}
