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

package printers_test

import (
	"encoding"
	"encoding/json"
	"fmt"

	"github.com/jtcressy/go-cli-toolkit/cmdutil/printers"
	"github.com/samber/lo"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FmtStringerStruct struct {
	FmtStringer
}

type FmtStringer func() string

var _ fmt.Stringer = (*FmtStringer)(nil)

func (s FmtStringer) String() string {
	return s()
}

type ErrorStruct struct {
	ExampleError
}

type ExampleError func() string

var _ error = (*ExampleError)(nil)
var _ error = (*ErrorStruct)(nil)

func (e ExampleError) Error() string {
	return e()
}

type TextMarshalerStruct struct {
	TextMarshaler
}

type TextMarshaler func() string

var _ encoding.TextMarshaler = (*TextMarshaler)(nil)

func (t TextMarshaler) MarshalText() ([]byte, error) {
	return []byte(t()), nil
}

type JSONMarshalerStruct struct {
	JSONMarshaler
}

type JSONMarshaler func() string

var _ json.Marshaler = (*JSONMarshaler)(nil)

func (j JSONMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(j()), nil
}

var _ = Describe("Table Reflection", Label("unit"), func() {
	type TestStruct struct {
		Field1 string `header:"Field 1"`
		Field2 int    `header:"Field 2"`
	}

	type TestCase struct {
		Input           interface{}
		ExpectedHeaders []string
		ExpectedRows    [][]string
		ShouldError     bool
	}

	DescribeTable(
		"GenerateTableData returns headers and rows for",
		func(tc TestCase) {
			headers, rows, err := printers.GenerateTableData(tc.Input)

			if tc.ShouldError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(headers).To(Equal(tc.ExpectedHeaders))
				Expect(rows).To(Equal(tc.ExpectedRows))
			}
		},
		// --- plain values
		// - string
		Entry("a string", TestCase{
			Input:           "hello",
			ExpectedHeaders: []string{},
			ExpectedRows:    [][]string{{"hello"}},
			ShouldError:     false,
		}),
		// - int
		Entry("an int", TestCase{
			Input:           42,
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"42"}},
			ShouldError:     false,
		}),
		// - bool
		Entry("a bool", TestCase{
			Input:           true,
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// - float64
		Entry("a float64", TestCase{
			Input:           3.14,
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"3.14"}},
			ShouldError:     false,
		}),
		// - complex128
		Entry("a complex128", TestCase{
			Input:           3.14i,
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"(0+3.14i)"}},
			ShouldError:     false,
		}),
		// - byte
		Entry("a byte", TestCase{
			Input:           byte('b'),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"98"}},
			ShouldError:     false,
		}),
		// - rune
		Entry("a rune", TestCase{
			Input:           'a',
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"97"}},
			ShouldError:     false,
		}),
		// - nil
		Entry("nil input", TestCase{
			Input:           nil,
			ExpectedHeaders: []string{},
			ExpectedRows:    [][]string{},
			ShouldError:     false,
		}),
		// - nil pointer
		Entry("a nil pointer", TestCase{
			Input:           (*string)(nil),
			ExpectedHeaders: []string{},
			ExpectedRows:    [][]string{},
			ShouldError:     false,
		}),
		// - nil interface
		Entry("a nil interface", TestCase{
			Input:           (any)(nil),
			ExpectedHeaders: []string{},
			ExpectedRows:    [][]string{},
			ShouldError:     false,
		}),
		// - nil double pointer
		Entry("a nil double pointer", TestCase{
			Input:           lo.ToPtr((*string)(nil)),
			ExpectedHeaders: []string{},
			ExpectedRows:    [][]string{},
			ShouldError:     false,
		}),
		// - nil pointer interface
		Entry("a nil pointer interface", TestCase{
			Input:           lo.ToPtr((any)(nil)),
			ExpectedHeaders: []string{},
			ExpectedRows:    [][]string{},
			ShouldError:     false,
		}),
		// --- basic struct
		Entry("a struct", TestCase{
			Input: TestStruct{
				Field1: "value1",
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// --- primitive types
		// - all primitive types as field values
		Entry("a struct with all primitive types", TestCase{
			Input: struct {
				Field1 string  `header:"Field 1"`
				Field2 int     `header:"Field 2"`
				Field3 bool    `header:"Field 3"`
				Field4 float64 `header:"Field 4"`
			}{
				Field1: "value1",
				Field2: 42,
				Field3: true,
				Field4: 3.14,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2", "Field 3", "Field 4"},
			ExpectedRows:    [][]string{{"value1", "42", "true", "3.14"}},
			ShouldError:     false,
		}),
		// - all primitive types as field pointers
		Entry("a struct with all primitive types as pointers", TestCase{
			Input: struct {
				Field1 *string  `header:"Field 1"`
				Field2 *int     `header:"Field 2"`
				Field3 *bool    `header:"Field 3"`
				Field4 *float64 `header:"Field 4"`
			}{
				Field1: lo.ToPtr("value1"),
				Field2: lo.ToPtr(42),
				Field3: lo.ToPtr(true),
				Field4: lo.ToPtr(3.14),
			},
			ExpectedHeaders: []string{"Field 1", "Field 2", "Field 3", "Field 4"},
			ExpectedRows:    [][]string{{"value1", "42", "true", "3.14"}},
			ShouldError:     false,
		}),
		// - all primitive types as nil field pointers
		Entry("a struct with all nil field pointers", TestCase{
			Input: struct {
				Field1 *string  `header:"Field 1"`
				Field2 *int     `header:"Field 2"`
				Field3 *bool    `header:"Field 3"`
				Field4 *float64 `header:"Field 4"`
			}{
				Field1: nil,
				Field2: nil,
				Field3: nil,
				Field4: nil,
			},
			ExpectedHeaders: nil,
			ExpectedRows:    [][]string{nil},
			ShouldError:     false,
		}),
		// --- primitive types in collections (slice/array)
		// - slice of primitive types
		Entry("a slice of primitive types", TestCase{
			Input: []any{
				true,              // bool
				42,                // int
				uint(42),          // uint
				float64(3.14),     // float64
				complex128(3.14i), // complex128
				"hello",           // string
				byte('b'),         // byte (alias for uint8)
				'a',               // rune (alias for int32)
			},
			ExpectedHeaders: []string{""},
			ExpectedRows: [][]string{
				{"true"},
				{"42"},
				{"42"},
				{"3.14"},
				{"(0+3.14i)"},
				{"hello"},
				{"98"},
				{"97"},
			},
			ShouldError: false,
		}),
		// - array of primitive types
		Entry("an array of primitive types", TestCase{
			Input: [8]any{
				true,              // bool
				42,                // int
				uint(42),          // uint
				float64(3.14),     // float64
				complex128(3.14i), // complex128
				"hello",           // string
				byte('b'),         // byte (alias for uint8)
				'a',               // rune (alias for int32)
			},
			ExpectedHeaders: []string{""},
			ExpectedRows: [][]string{
				{"true"},
				{"42"},
				{"42"},
				{"3.14"},
				{"(0+3.14i)"},
				{"hello"},
				{"98"},
				{"97"},
			},
			ShouldError: false,
		}),
		// - slice of primitive type pointers
		Entry("a slice of primitive type pointers", TestCase{
			Input: []any{
				lo.ToPtr(true),
				lo.ToPtr(42),
				lo.ToPtr(uint(42)),
				lo.ToPtr(3.14),
				lo.ToPtr(complex128(3.14i)),
				lo.ToPtr("hello"),
				lo.ToPtr(byte('b')),
				lo.ToPtr('a'),
			},
			ExpectedHeaders: []string{""},
			ExpectedRows: [][]string{
				{"true"},
				{"42"},
				{"42"},
				{"3.14"},
				{"(0+3.14i)"},
				{"hello"},
				{"98"},
				{"97"},
			},
			ShouldError: false,
		}),
		// - array of primitive type pointers
		Entry("an array of primitive type pointers", TestCase{
			Input: [8]any{
				lo.ToPtr(true),
				lo.ToPtr(42),
				lo.ToPtr(uint(42)),
				lo.ToPtr(3.14),
				lo.ToPtr(complex128(3.14i)),
				lo.ToPtr("hello"),
				lo.ToPtr(byte('b')),
				lo.ToPtr('a'),
			},
			ExpectedHeaders: []string{""},
			ExpectedRows: [][]string{
				{"true"},
				{"42"},
				{"42"},
				{"3.14"},
				{"(0+3.14i)"},
				{"hello"},
				{"98"},
				{"97"},
			},
			ShouldError: false,
		}),
		// - slice of nil primitive type pointers
		Entry("a slice of nil primitive type pointers", TestCase{
			Input: []any{
				(*bool)(nil),
				(*int)(nil),
				(*uint)(nil),
				(*float64)(nil),
				(*complex128)(nil),
				(*string)(nil),
				(*byte)(nil),
				(*rune)(nil),
			},
			ExpectedHeaders: []string{""},
			ExpectedRows:    [][]string{},
			ShouldError:     false,
		}),
		// - array of nil primitive type pointers
		Entry("an array of nil primitive type pointers", TestCase{
			Input: [8]any{
				(*bool)(nil),
				(*int)(nil),
				(*uint)(nil),
				(*float64)(nil),
				(*complex128)(nil),
				(*string)(nil),
				(*byte)(nil),
				(*rune)(nil),
			},
			ExpectedHeaders: []string{""},
			ExpectedRows:    [][]string{},
			ShouldError:     false,
		}),
		// --- structs in collections (slice/array)
		// - slice of structs
		Entry("a slice of structs", TestCase{
			Input: []TestStruct{
				{
					Field1: "value1",
					Field2: 42,
				},
				{
					Field1: "value2",
					Field2: 43,
				},
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}, {"value2", "43"}},
			ShouldError:     false,
		}),
		// - array of structs
		Entry("an array of structs", TestCase{
			Input: [2]TestStruct{
				{
					Field1: "value1",
					Field2: 42,
				},
				{
					Field1: "value2",
					Field2: 43,
				},
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}, {"value2", "43"}},
			ShouldError:     false,
		}),
		// - slice of pointer structs
		Entry("a slice of pointer structs", TestCase{
			Input: []*TestStruct{
				{
					Field1: "value1",
					Field2: 42,
				},
				{
					Field1: "value2",
					Field2: 43,
				},
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}, {"value2", "43"}},
			ShouldError:     false,
		}),
		// - array of pointer structs
		Entry("an array of pointer structs", TestCase{
			Input: [2]*TestStruct{
				{
					Field1: "value1",
					Field2: 42,
				},
				{
					Field1: "value2",
					Field2: 43,
				},
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}, {"value2", "43"}},
			ShouldError:     false,
		}),
		// --- structs with embedded structs
		// - struct with embedded struct field
		Entry("a struct with an embedded struct field", TestCase{
			Input: struct {
				TestStruct `header:"TestStruct"`
				Field3     bool `header:"Field 3"`
			}{
				TestStruct: TestStruct{
					Field1: "value1",
					Field2: 42,
				},
				Field3: true,
			},
			ExpectedHeaders: []string{"TestStruct Field 1", "TestStruct Field 2", "Field 3"},
			ExpectedRows:    [][]string{{"value1", "42", "true"}},
			ShouldError:     false,
		}),
		// - struct with embedded pointer struct field
		Entry("a struct with an embedded pointer struct field", TestCase{
			Input: struct {
				*TestStruct `header:"TestStruct"`
				Field3      bool `header:"Field 3"`
			}{
				TestStruct: &TestStruct{
					Field1: "value1",
					Field2: 42,
				},
				Field3: true,
			},
			ExpectedHeaders: []string{"TestStruct Field 1", "TestStruct Field 2", "Field 3"},
			ExpectedRows:    [][]string{{"value1", "42", "true"}},
			ShouldError:     false,
		}),
		// - struct with embedded pointer struct field with nil value
		Entry("a struct with an embedded pointer struct field with nil value", TestCase{
			Input: struct {
				*TestStruct `header:"TestStruct"`
				Field3      bool `header:"Field 3"`
			}{
				TestStruct: nil,
				Field3:     true,
			},
			ExpectedHeaders: []string{"Field 3"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// --- nested structs (structs with named struct fields)
		// - struct with nested struct field
		Entry("a struct with a nested struct field", TestCase{
			Input: struct {
				TestStruct struct {
					Field1 string `header:"Field 1"`
					Field2 int    `header:"Field 2"`
				} `header:"TestStruct"`
				Field3 bool `header:"Field 3"`
			}{
				TestStruct: struct {
					Field1 string `header:"Field 1"`
					Field2 int    `header:"Field 2"`
				}{
					Field1: "value1",
					Field2: 42,
				},
				Field3: true,
			},
			ExpectedHeaders: []string{"Field 3"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// - struct with nested pointer struct field
		Entry("a struct with a nested pointer struct field", TestCase{
			Input: struct {
				TestStruct *struct {
					Field1 string `header:"Field 1"`
					Field2 int    `header:"Field 2"`
				} `header:"TestStruct"`
				Field3 bool `header:"Field 3"`
			}{
				TestStruct: &struct {
					Field1 string `header:"Field 1"`
					Field2 int    `header:"Field 2"`
				}{
					Field1: "value1",
					Field2: 42,
				},
				Field3: true,
			},
			ExpectedHeaders: []string{"Field 3"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// - struct with nested pointer struct field with nil value
		Entry("a struct with a nested pointer struct field with nil value", TestCase{
			Input: struct {
				TestStruct *struct {
					Field1 string `header:"Field 1"`
					Field2 int    `header:"Field 2"`
				} `header:"TestStruct"`
				Field3 bool `header:"Field 3"`
			}{
				TestStruct: nil,
				Field3:     true,
			},
			ExpectedHeaders: []string{"Field 3"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// --- validate header formatting (no tag, header tag and json tag)
		// - no tag
		Entry("a struct with no tag", TestCase{
			Input: struct {
				Field1 string
				Field2 int
			}{
				Field1: "value1",
				Field2: 42,
			},
			ExpectedHeaders: []string{"FIELD 1", "FIELD 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - header tag
		Entry("a struct with header tag", TestCase{
			Input: struct {
				Field1 string `header:"Field 1"`
				Field2 int    `header:"Field 2"`
			}{
				Field1: "value1",
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - json tag
		Entry("a struct with json tag", TestCase{
			Input: struct {
				Field1 string `json:"field1"`
				Field2 int    `json:"field2"`
			}{
				Field1: "value1",
				Field2: 42,
			},
			ExpectedHeaders: []string{"FIELD 1", "FIELD 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// --- mixed data structures
		// - struct with nested slice of structs
		Entry("a struct with a nested slice of structs", TestCase{
			Input: struct {
				TestStructs []TestStruct `header:"TestStructs"`
				Field3      bool         `header:"Field 3"`
			}{
				TestStructs: []TestStruct{
					{
						Field1: "value1",
						Field2: 42,
					},
					{
						Field1: "value2",
						Field2: 43,
					},
				},
				Field3: true,
			},
			ExpectedHeaders: []string{"TestStructs", "Field 3"},
			ExpectedRows:    [][]string{{"[{value1 42} {value2 43}]", "true"}},
			ShouldError:     false,
		}),
		// --- fields that implement supported interfaces
		// - as single struct
		Entry("a struct with fields that implement supported interfaces", TestCase{
			Input: struct {
				Field1 FmtStringer    `header:"Field 1"`
				Field2 *FmtStringer   `header:"Field 2"`
				Field3 ExampleError   `header:"Field 3"`
				Field4 *ExampleError  `header:"Field 4"`
				Field5 TextMarshaler  `header:"Field 5"`
				Field6 *TextMarshaler `header:"Field 6"`
				Field7 JSONMarshaler  `header:"Field 7"`
				Field8 *JSONMarshaler `header:"Field 8"`
			}{
				Field1: FmtStringer(func() string { return "value1" }),
				Field2: lo.ToPtr(FmtStringer(func() string { return "value2" })),
				Field3: ExampleError(func() string { return "error" }),
				Field4: lo.ToPtr(ExampleError(func() string { return "error" })),
				Field5: TextMarshaler(func() string { return "value5" }),
				Field6: lo.ToPtr(TextMarshaler(func() string { return "value6" })),
				Field7: JSONMarshaler(func() string { return "value7" }),
				Field8: lo.ToPtr(JSONMarshaler(func() string { return "value8" })),
			},
			ExpectedHeaders: []string{
				"Field 1",
				"Field 2",
				"Field 3",
				"Field 4",
				"Field 5",
				"Field 6",
				"Field 7",
				"Field 8",
			},
			ExpectedRows: [][]string{
				{"value1", "value2", "error", "error", "value5", "value6", "value7", "value8"},
			},
			ShouldError: false,
		}),
		// - as struct with fields that implement supported interfaces using struct as underlying
		// type
		Entry(
			"a struct with fields that implement supported interfaces using struct as underlying type",
			TestCase{
				Input: struct {
					Field1 FmtStringerStruct    `header:"Field 1"`
					Field2 *FmtStringerStruct   `header:"Field 2"`
					Field3 ErrorStruct          `header:"Field 3"`
					Field4 *ErrorStruct         `header:"Field 4"`
					Field5 TextMarshalerStruct  `header:"Field 5"`
					Field6 *TextMarshalerStruct `header:"Field 6"`
					Field7 JSONMarshalerStruct  `header:"Field 7"`
					Field8 *JSONMarshalerStruct `header:"Field 8"`
				}{
					Field1: FmtStringerStruct{FmtStringer(func() string { return "value1" })},
					Field2: lo.ToPtr(
						FmtStringerStruct{FmtStringer(func() string { return "value2" })},
					),
					Field3: ErrorStruct{ExampleError(func() string { return "error" })},
					Field4: lo.ToPtr(ErrorStruct{ExampleError(func() string { return "error" })}),
					Field5: TextMarshalerStruct{TextMarshaler(func() string { return "value5" })},
					Field6: lo.ToPtr(
						TextMarshalerStruct{TextMarshaler(func() string { return "value6" })},
					),
					Field7: JSONMarshalerStruct{JSONMarshaler(func() string { return "value7" })},
					Field8: lo.ToPtr(
						JSONMarshalerStruct{JSONMarshaler(func() string { return "value8" })},
					),
				},
				ExpectedHeaders: []string{
					"Field 1",
					"Field 2",
					"Field 3",
					"Field 4",
					"Field 5",
					"Field 6",
					"Field 7",
					"Field 8",
				},
				ExpectedRows: [][]string{
					{"value1", "value2", "error", "error", "value5", "value6", "value7", "value8"},
				},
				ShouldError: false,
			},
		),
		// - as slice
		Entry("a slice of objects that implement supported interfaces", TestCase{
			Input: []any{
				FmtStringer(func() string { return "value1" }),
				lo.ToPtr(FmtStringer(func() string { return "value2" })),
				ExampleError(func() string { return "error" }),
				lo.ToPtr(ExampleError(func() string { return "error" })),
				TextMarshaler(func() string { return "value5" }),
				lo.ToPtr(TextMarshaler(func() string { return "value6" })),
				JSONMarshaler(func() string { return "value7" }),
				lo.ToPtr(JSONMarshaler(func() string { return "value8" })),
			},
			ExpectedHeaders: []string{""},
			ExpectedRows: [][]string{
				{"value1"},
				{"value2"},
				{"error"},
				{"error"},
				{"value5"},
				{"value6"},
				{"value7"},
				{"value8"},
			},
			ShouldError: false,
		}),
		// -- fmt.Stringer
		// - as field value
		Entry("a field that implements fmt.Stringer", TestCase{
			Input: struct {
				Field1 FmtStringer `header:"Field 1"`
				Field2 int         `header:"Field 2"`
			}{
				Field1: func() string { return "value1" },
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field pointer
		Entry("a field pointer that implements fmt.Stringer", TestCase{
			Input: struct {
				Field1 *FmtStringer `header:"Field 1"`
				Field2 int          `header:"Field 2"`
			}{
				Field1: lo.ToPtr(FmtStringer(func() string { return "value1" })),
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field double pointer
		Entry("a field double pointer that implements fmt.Stringer", TestCase{
			Input: struct {
				Field1 **any `header:"Field 1"`
				Field2 int   `header:"Field 2"`
			}{
				Field1: lo.ToPtr(lo.ToPtr(any(FmtStringer(func() string { return "value1" })))),
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field struct
		Entry("a struct that implements fmt.Stringer", TestCase{
			Input: struct {
				Field1 FmtStringerStruct `header:"Field 1"`
				Field2 int               `header:"Field 2"`
			}{
				Field1: FmtStringerStruct{FmtStringer(func() string { return "value1" })},
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as single value
		Entry("a single value that implements fmt.Stringer", TestCase{
			Input:           FmtStringer(func() string { return "value1" }),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"value1"}},
			ShouldError:     false,
		}),
		// - as pointer
		Entry("a pointer that implements fmt.Stringer", TestCase{
			Input:           lo.ToPtr(FmtStringer(func() string { return "value1" })),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"value1"}},
			ShouldError:     false,
		}),
		// -- error
		// - as field value
		Entry("a field that implements error", TestCase{
			Input: struct {
				Field1 ExampleError `header:"Field 1"`
				Field2 int          `header:"Field 2"`
			}{
				Field1: func() string { return "error" },
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"error", "42"}},
			ShouldError:     false,
		}),
		// - as field pointer
		Entry("a field pointer that implements error", TestCase{
			Input: struct {
				Field1 *ExampleError `header:"Field 1"`
				Field2 int           `header:"Field 2"`
			}{
				Field1: lo.ToPtr(ExampleError(func() string { return "error" })),
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"error", "42"}},
			ShouldError:     false,
		}),
		// - as field double pointer
		Entry("a field double pointer that implements error", TestCase{
			Input: struct {
				Field1 **any `header:"Field 1"`
				Field2 int   `header:"Field 2"`
			}{
				Field1: lo.ToPtr(lo.ToPtr(any(ExampleError(func() string { return "error" })))),
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"error", "42"}},
			ShouldError:     false,
		}),
		// - as field struct
		Entry("a struct that implements error", TestCase{
			Input: struct {
				Field1 ErrorStruct `header:"Field 1"`
				Field2 int         `header:"Field 2"`
			}{
				Field1: ErrorStruct{ExampleError(func() string { return "error" })},
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"error", "42"}},
			ShouldError:     false,
		}),
		// - as single value
		Entry("a single value that implements error", TestCase{
			Input:           ExampleError(func() string { return "error" }),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"error"}},
			ShouldError:     false,
		}),
		// - as pointer
		Entry("a pointer that implements error", TestCase{
			Input:           lo.ToPtr(ExampleError(func() string { return "error" })),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"error"}},
			ShouldError:     false,
		}),
		// -- encoding.TextMarshaler
		// - as field value
		Entry("a field that implements encoding.TextMarshaler", TestCase{
			Input: struct {
				Field1 TextMarshaler `header:"Field 1"`
				Field2 int           `header:"Field 2"`
			}{
				Field1: func() string { return "value1" },

				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field pointer
		Entry("a field pointer that implements encoding.TextMarshaler", TestCase{
			Input: struct {
				Field1 *TextMarshaler `header:"Field 1"`
				Field2 int            `header:"Field 2"`
			}{
				Field1: lo.ToPtr(TextMarshaler(func() string { return "value1" })),
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field double pointer
		Entry("a field double pointer that implements encoding.TextMarshaler", TestCase{
			Input: struct {
				Field1 **any `header:"Field 1"`
				Field2 int   `header:"Field 2"`
			}{
				Field1: lo.ToPtr(lo.ToPtr(any(TextMarshaler(func() string { return "value1" })))),
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field struct
		Entry("a struct that implements encoding.TextMarshaler", TestCase{
			Input: struct {
				Field1 TextMarshalerStruct `header:"Field 1"`
				Field2 int                 `header:"Field 2"`
			}{
				Field1: TextMarshalerStruct{TextMarshaler(func() string { return "value1" })},
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as single value

		Entry("a single value that implements encoding.TextMarshaler", TestCase{
			Input:           TextMarshaler(func() string { return "value1" }),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"value1"}},
			ShouldError:     false,
		}),
		// - as pointer
		Entry("a pointer that implements encoding.TextMarshaler", TestCase{
			Input:           lo.ToPtr(TextMarshaler(func() string { return "value1" })),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"value1"}},
			ShouldError:     false,
		}),
		// -- encoding/json.Marshaler
		// - as field value
		Entry("a field that implements encoding/json.Marshaler", TestCase{
			Input: struct {
				Field1 JSONMarshaler `header:"Field 1"`
				Field2 int           `header:"Field 2"`
			}{
				Field1: func() string { return "value1" },
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field pointer
		Entry("a field pointer that implements encoding/json.Marshaler", TestCase{
			Input: struct {
				Field1 *JSONMarshaler `header:"Field 1"`
				Field2 int            `header:"Field 2"`
			}{
				Field1: lo.ToPtr(JSONMarshaler(func() string { return "value1" })),
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field double pointer
		Entry("a field double pointer that implements encoding/json.Marshaler", TestCase{
			Input: struct {
				Field1 **any `header:"Field 1"`
				Field2 int   `header:"Field 2"`
			}{
				Field1: lo.ToPtr(lo.ToPtr(any(JSONMarshaler(func() string { return "value1" })))),
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as field struct
		Entry("a struct that implements encoding/json.Marshaler", TestCase{
			Input: struct {
				Field1 JSONMarshalerStruct `header:"Field 1"`
				Field2 int                 `header:"Field 2"`
			}{
				Field1: JSONMarshalerStruct{JSONMarshaler(func() string { return "value1" })},
				Field2: 42,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - as single value
		Entry("a single value that implements encoding/json.Marshaler", TestCase{
			Input:           JSONMarshaler(func() string { return "value1" }),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"value1"}},
			ShouldError:     false,
		}),
		// - as pointer
		Entry("a pointer that implements encoding/json.Marshaler", TestCase{
			Input:           lo.ToPtr(JSONMarshaler(func() string { return "value1" })),
			ExpectedHeaders: []string{"Value"},
			ExpectedRows:    [][]string{{"value1"}},
			ShouldError:     false,
		}),
		// --- extra tag features
		// -- inline header tag
		// - embedded struct
		Entry("an embedded struct with inline header tag", TestCase{
			Input: struct {
				TestStruct `header:"Test Struct,inline"`
				Field3     bool `header:"Field 3"`
			}{
				TestStruct: TestStruct{
					Field1: "value1",
					Field2: 42,
				},
				Field3: true,
			},
			ExpectedHeaders: []string{"Test Struct Field 1", "Test Struct Field 2", "Field 3"},
			ExpectedRows:    [][]string{{"value1", "42", "true"}},
			ShouldError:     false,
		}),
		// - named struct
		Entry("a named struct with inline header tag", TestCase{
			Input: struct {
				Field1 TestStruct `header:"Test,inline"`
				Field2 bool       `header:"Field 2"`
			}{
				Field1: TestStruct{
					Field1: "value1",
					Field2: 42,
				},
				Field2: true,
			},
			ExpectedHeaders: []string{"Test Field 1", "Test Field 2", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42", "true"}},
			ShouldError:     false,
		}),
		// -- omit tag
		// - embedded struct
		Entry("an embedded struct with omit tag", TestCase{
			Input: struct {
				TestStruct `header:"-,omit"`
				Field3     bool `header:"Field 3"`
			}{
				TestStruct: TestStruct{
					Field1: "value1",
					Field2: 42,
				},
				Field3: true,
			},
			ExpectedHeaders: []string{"Field 3"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// - named struct
		Entry("a named struct with omit tag", TestCase{
			Input: struct {
				Field1 TestStruct `header:"-"`
				Field2 bool       `header:"Field 2"`
			}{
				Field1: TestStruct{
					Field1: "value1",
					Field2: 42,
				},
				Field2: true,
			},
			ExpectedHeaders: []string{"Field 2"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// -- json omit tag
		// - embedded struct
		Entry("an embedded struct with json omit tag", TestCase{
			Input: struct {
				TestStruct `json:"-"`
				Field3     bool `header:"Field 3"`
			}{
				TestStruct: TestStruct{
					Field1: "value1",
					Field2: 42,
				},
				Field3: true,
			},
			ExpectedHeaders: []string{"Field 3"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// - named struct
		Entry("a named struct with json omit tag", TestCase{
			Input: struct {
				Field1 TestStruct `json:"-"`
				Field2 bool       `header:"Field 2"`
			}{
				Field1: TestStruct{
					Field1: "value1",
					Field2: 42,
				},
				Field2: true,
			},
			ExpectedHeaders: []string{"Field 2"},
			ExpectedRows:    [][]string{{"true"}},
			ShouldError:     false,
		}),
		// --- edge cases
		// - double pointer interface
		Entry("a double pointer interface", TestCase{
			Input:           lo.ToPtr(lo.ToPtr(any("hello"))),
			ExpectedHeaders: []string{},
			ExpectedRows:    [][]string{{"hello"}},
			ShouldError:     false,
		}),
		// - pointer struct slice with nil element
		Entry("a pointer struct slice with nil element", TestCase{
			Input: []*TestStruct{
				{
					Field1: "value1",
					Field2: 42,
				},
				nil,
			},
			ExpectedHeaders: []string{"Field 1", "Field 2"},
			ExpectedRows:    [][]string{{"value1", "42"}},
			ShouldError:     false,
		}),
		// - nil pointer struct
		Entry("a nil pointer struct", TestCase{
			Input:           (*TestStruct)(nil),
			ExpectedHeaders: []string{},
			ExpectedRows:    [][]string{},
			ShouldError:     false,
		}),
	)
})
