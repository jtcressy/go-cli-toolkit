// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Joel Cressy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package printers

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

var (
	DefaultTableReflectorFunc TableReflectorFunc = GenerateTableData
)

// TableReflectorFunc should take any struct or collection of structs and use
// reflection to resolve tabular data. Depending on the implementation, this may
// use struct field tags similar to encoding/json to declare table header labels.
type TableReflectorFunc func(data any) (headers []string, rows [][]string, err error)

// GenerateTableData is a TableReflectorFunc that will use reflection to resolve
// tabular data from a struct tor slice/array of structs.
// Every field in the input struct represents a column on the table.
//
// # Defining Headers
//
// Input struct(s) should use field tags to define table headers at compile-time.
//
// Note: json field tags are parsed from various snake-case/camel-case/title-case/etc.
// to normalized all-caps (shouting) strings that isolate the words for easier reading and viewing.
// Header tags are not subject to this normalization, and can be considered overrides of this
// behavior.
//
// There are 3 main ways headers are resolved:
//
// 1. header field tags
//
//	MyField `header:"MY FIELD"` -> "MY FIELD"
//	MyField `header:"something completely different!"` -> "something completely different!"
//
// 2. json field tags
//
//	MyField `json:"my_field"` -> "MY FIELD"
//	MyField `json:"myField"` -> "MY FIELD"
//
// 3. The field name itself
//
//	MyField -> "MY FIELD"
//
// Complex data structures or nested/embedded structs:
//
//	type MyStruct struct {
//	    Field1 string `header:"Field 1"`
//	    Field2 int    `header:"Field 2"`
//	}
//
//	type MyOtherStruct struct {
//	    MyStruct
//	    Field3 bool `header:"Field 3"`
//	}
//
//	myData := &MyOtherStruct{
//	    MyStruct: MyStruct{
//	        Field1: "value1",
//	        Field2: 42,
//	    },
//	    Field3: true,
//	}
//
//	headers, rows, err := GenerateTableData(myData)
//
//	// headers == []string{"Field 1", "Field 2", "Field 3"}
//	// rows == [][]string{{"value1", "42", "true"}}
//
//	// example csv output:
//	// Field 1,Field 2,Field 3
//	// value1,42,true
func GenerateTableData(data any) (headers []string, rows [][]string, _ error) {
	headers = []string{}
	rows = [][]string{}
	if data == nil {
		return headers, rows, nil
	}
	if v, ok := indirectValue(reflect.ValueOf(data)); ok {
		headers, rows = processValue(v)
	}
	return headers, rows, nil
}

func processValue(v reflect.Value) (headers []string, rows [][]string) {
	headers = []string{}
	rows = [][]string{}
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		headers = []string{""}
		for i := 0; i < v.Len(); i++ {
			if fv, ok := indirectValue(v.Index(i)); ok {
				switch fv.Kind() {
				case reflect.String:
					rows = append(rows, []string{reflect.Indirect(fv).String()})
				case reflect.Struct:
					itemHeaders, itemRow := recursiveFieldExtract(fv, "", false)
					if i == 0 {
						headers = itemHeaders
					}
					rows = append(rows, itemRow)
				case reflect.Array, reflect.Slice:
					// TODO: someday support nested tables, like html tables
					fallthrough
				case reflect.Invalid,
					reflect.Bool,
					reflect.Int,
					reflect.Int8,
					reflect.Int16,
					reflect.Int32,
					reflect.Int64,
					reflect.Uint,
					reflect.Uint8,
					reflect.Uint16,
					reflect.Uint32,
					reflect.Uint64,
					reflect.Uintptr,
					reflect.Float32,
					reflect.Float64,
					reflect.Complex64,
					reflect.Complex128,
					reflect.Chan,
					reflect.Func,
					reflect.Interface,
					reflect.Map,
					reflect.Ptr,
					reflect.UnsafePointer:
					fallthrough
				default:
					rows = append(rows, []string{resolveStringValue(fv)})
				}
			}
		}
	case reflect.String:
		rows = [][]string{{v.String()}}
	case reflect.Struct:
		h, r := recursiveFieldExtract(v, "", false)
		headers = h
		rows = append(rows, r)
	case reflect.Invalid, reflect.Chan, reflect.Ptr, reflect.UnsafePointer:
		break
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Interface, reflect.Map, reflect.Func:
		fallthrough
	default:
		val := resolveStringValue(v)
		if len(val) > 0 {
			rows = [][]string{{val}}
			headers = []string{"Value"}
		}
	}
	return headers, rows
}

func resolveStringValue(v reflect.Value) string {
	if val, ok := resolveStringerInterfaces(v); ok {
		return val
	} else if v.CanInterface() && v.Kind() != reflect.Func {
		return fmt.Sprint(v.Interface())
	}
	return ""
}

// indirectValue returns the indirect value of a [reflect.Value], handling pointers and interfaces
// (including nested, like double pointers or interfaces of interfaces)
// Returns false if the value is nil.
// Similar to [reflect.Indirect], but it also handles interfaces, and it resolves nested
// pointers/interfaces of any depth.
func indirectValue(v reflect.Value) (val reflect.Value, isNotNil bool) {
	val = v
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return val, false
		}
		val = val.Elem()
	}
	return val, true
}

func recursiveFieldExtract(
	v reflect.Value,
	prefix string,
	inline bool,
) (headers, row []string) {
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := typ.Field(i)
		if fv, ok := indirectValue(v.Field(i)); ok {
			header := resolveHeader(f, prefix, inline)
			// If header is marked with '-', skip this field
			if header == "-" {
				continue
			}
			headers, row = extractField(fv, f, header, headers, row)
		}
	}
	return headers, row
}

func resolveHeader(f reflect.StructField, prefix string, inline bool) string {
	var header string
	tag := f.Tag.Get("header")
	if tag != "" {
		header = strings.Split(tag, ",")[0]
	} else if tag = f.Tag.Get("json"); tag != "" {
		header = strings.Split(tag, ",")[0]
		if header == "-" {
			return "-"
		}
		header = strcase.ToScreamingDelimited(header, ' ', "", true)
	}
	if header == "" {
		header = strcase.ToScreamingDelimited(f.Name, ' ', "", true)
	}
	if inline {
		header = prefix + " " + header
	}
	return header
}

func extractField(
	fv reflect.Value,
	f reflect.StructField,
	header string,
	headers, row []string,
) ([]string, []string) {
	switch fv.Kind() {
	case reflect.String:
		headers = append(headers, header)
		row = append(row, reflect.Indirect(fv).String())
	case reflect.Struct:
		if f.Anonymous || strings.Contains(f.Tag.Get("header"), ",inline") {
			embHeaders, embRow := recursiveFieldExtract(fv, header, true)
			headers = append(headers, embHeaders...)
			row = append(row, embRow...)
		} else if val, ok := resolveStringerInterfaces(fv); ok {
			headers = append(headers, header)
			row = append(row, val)
		}
	case reflect.Array, reflect.Slice:
		// TODO: someday support nested tables, like html tables
		fallthrough
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		fallthrough
	default:
		if val := resolveStringValue(fv); val != "" {
			headers = append(headers, header)
			row = append(row, val)
		}
	}
	return headers, row
}

func resolveStringerInterfaces(v reflect.Value) (string, bool) {
	if val, ok := reflectCast[fmt.Stringer](v, v.CanAddr()); ok {
		return val.String(), true
	}
	if val, ok := reflectCast[error](v, v.CanAddr()); ok {
		return val.Error(), true
	}
	if val, ok := reflectCast[encoding.TextMarshaler](v, v.CanAddr()); ok {
		if bts, err := val.MarshalText(); err == nil {
			return string(bts), true
		}
	}
	if val, ok := reflectCast[json.Marshaler](v, v.CanAddr()); ok {
		if bts, err := val.MarshalJSON(); err == nil {
			return string(bts), true
		}
	}
	return "", false
}

func reflectCast[T any](v reflect.Value, allowAddr bool) (T, bool) {
	t := v.Type()
	it := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(it) {
		return reflectCastFromAddr[T](v)
	}
	if t.Implements(it) {
		return reflectCastFromValue[T](v)
	}
	return *new(T), false
}

func reflectCastFromAddr[T any](v reflect.Value) (T, bool) {
	va := v.Addr()
	m, ok := va.Interface().(T)
	return m, ok
}

func reflectCastFromValue[T any](v reflect.Value) (T, bool) {
	m, ok := v.Interface().(T)
	return m, ok
}
