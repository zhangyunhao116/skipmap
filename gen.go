//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	_ "embed"
)

// Inspired by sort/gen_sort_variants.go
type Variant struct {
	// Package is the package name.
	Package string

	// Name is the variant name: should be unique among variants.
	Name string

	// Path is the file path into which the generator will emit the code for this
	// variant.
	Path string

	// Imports is the imports needed for this package.
	Imports string

	StructPrefix    string
	StructPrefixLow string
	StructSuffix    string
	ExtraFileds     string

	// Basic key and value type.
	KeyType   string
	ValueType string

	// Basic type argument.
	TypeArgument string

	// TypeParam is the optional type parameter for the function.
	TypeParam string // e.g. [T any]

	// Funcs is a map of functions used from within the template. The following
	// functions are expected to exist:
	Funcs template.FuncMap
}

type TypeReplacement struct {
	Type string
	Desc string
}

func main() {
	// For New.
	base := &Variant{
		Package:         "skipmap",
		Name:            "ordered",
		Path:            "gen_ordered.go",
		Imports:         "\"sync\"\n\"sync/atomic\"\n\"unsafe\"\n",
		KeyType:         "keyT",
		ValueType:       "valueT",
		TypeArgument:    "[keyT, valueT]",
		TypeParam:       "[keyT ordered, valueT any]",
		StructPrefix:    "Ordered",
		StructPrefixLow: "ordered",
		StructSuffix:    "",
		Funcs: template.FuncMap{
			"Less": func(i, j string) string {
				return fmt.Sprintf("(%s < %s)", i, j)
			},
			"Equal": func(i, j string) string {
				return fmt.Sprintf("%s == %s", i, j)
			},
		},
	}
	generate(base)
	base.Name += "Desc"
	base.StructSuffix += "Desc"
	base.Path = "gen_ordereddesc.go"
	base.Funcs = template.FuncMap{
		"Less": func(i, j string) string {
			return fmt.Sprintf("(%s > %s)", i, j)
		},
		"Equal": func(i, j string) string {
			return fmt.Sprintf("%s == %s", i, j)
		},
	}
	generate(base)

	// For NewFunc.
	basefunc := &Variant{
		Package:         "skipmap",
		Name:            "func",
		Path:            "gen_func.go",
		Imports:         "\"sync\"\n\"sync/atomic\"\n\"unsafe\"\n",
		KeyType:         "keyT",
		ValueType:       "valueT",
		TypeArgument:    "[keyT, valueT]",
		TypeParam:       "[keyT ordered, valueT any]",
		ExtraFileds:     "\nless func(a,b keyT)bool\n",
		StructPrefix:    "Func",
		StructPrefixLow: "func",
		StructSuffix:    "",
		Funcs: template.FuncMap{
			"Less": func(i, j string) string {
				return fmt.Sprintf("s.less(%s,%s)", i, j)
			},
			"Equal": func(i, j string) string {
				return fmt.Sprintf("!s.less(%s,%s)", j, i)
			},
		},
	}
	generate(basefunc)

	// For New{{Type}}.
	ts := []string{"String", "Float32", "Float64", "Int", "Int64", "Int32", "Uint64", "Uint32", "Uint"}
	for _, t := range ts {
		baseType := &Variant{
			Package:         "skipmap",
			Name:            "{{TypeLow}}",
			Path:            "gen_{{TypeLow}}.go",
			Imports:         "\"sync\"\n\"sync/atomic\"\n\"unsafe\"\n",
			KeyType:         "{{TypeLow}}",
			ValueType:       "valueT",
			TypeArgument:    "[valueT]",
			TypeParam:       "[valueT any]",
			StructPrefix:    "{{Type}}",
			StructPrefixLow: "{{TypeLow}}",
			StructSuffix:    "",
			Funcs: template.FuncMap{
				"Less": func(i, j string) string {
					return fmt.Sprintf("(%s < %s)", i, j)
				},
				"Equal": func(i, j string) string {
					return fmt.Sprintf("%s == %s", i, j)
				},
			},
		}
		baseTypeDesc := &Variant{
			Package:         "skipmap",
			Name:            "{{TypeLow}}Desc",
			Path:            "gen_{{TypeLow}}desc.go",
			Imports:         "\"sync\"\n\"sync/atomic\"\n\"unsafe\"\n",
			KeyType:         "{{TypeLow}}",
			ValueType:       "valueT",
			TypeArgument:    "[valueT]",
			TypeParam:       "[valueT any]",
			StructPrefix:    "{{Type}}",
			StructPrefixLow: "{{TypeLow}}",
			StructSuffix:    "Desc",
			Funcs: template.FuncMap{
				"Less": func(i, j string) string {
					return fmt.Sprintf("(%s > %s)", i, j)
				},
				"Equal": func(i, j string) string {
					return fmt.Sprintf("%s == %s", i, j)
				},
			},
		}
		tl := strings.ToLower(t)
		baseType.StructPrefix = strings.Replace(baseType.StructPrefix, "{{Type}}", t, -1)
		baseType.Name = strings.Replace(baseType.Name, "{{TypeLow}}", tl, -1)
		baseType.Path = strings.Replace(baseType.Path, "{{TypeLow}}", tl, -1)
		baseType.KeyType = strings.Replace(baseType.KeyType, "{{TypeLow}}", tl, -1)
		baseType.StructPrefixLow = strings.Replace(baseType.StructPrefixLow, "{{TypeLow}}", tl, -1)

		baseTypeDesc.StructPrefix = strings.Replace(baseTypeDesc.StructPrefix, "{{Type}}", t, -1)
		baseTypeDesc.Name = strings.Replace(baseTypeDesc.Name, "{{TypeLow}}", tl, -1)
		baseTypeDesc.Path = strings.Replace(baseTypeDesc.Path, "{{TypeLow}}", tl, -1)
		baseTypeDesc.KeyType = strings.Replace(baseTypeDesc.KeyType, "{{TypeLow}}", tl, -1)
		baseTypeDesc.StructPrefixLow = strings.Replace(baseTypeDesc.StructPrefixLow, "{{TypeLow}}", tl, -1)

		generate(baseType)
		generate(baseTypeDesc)
	}
}

// generate generates the code for variant `v` into a file named by `v.Path`.
func generate(v *Variant) {
	// Parse templateCode anew for each variant because Parse requires Funcs to be
	// registered, and it helps type-check the funcs.
	tmpl, err := template.New("gen").Funcs(v.Funcs).Parse(templateCode)
	if err != nil {
		log.Fatal("template Parse:", err)
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, v)
	if err != nil {
		log.Fatal("template Execute:", err)
	}

	os.WriteFile(v.Path, out.Bytes(), 0644)

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		println(string(out.Bytes()))
		log.Fatal("format:", err)
	}

	if err := os.WriteFile(v.Path, formatted, 0644); err != nil {
		log.Fatal("WriteFile:", err)
	}
}

//go:embed skipmap.tpl
var templateCode string
