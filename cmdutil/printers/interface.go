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

import "io"

// ObjectPrinterFunc is a function that can print objects.
type ObjectPrinterFunc func(obj any, w io.Writer) error

// PrintObj implements ObjectPrinter.
func (fn ObjectPrinterFunc) PrintObj(obj any, w io.Writer) error {
	return fn(obj, w)
}

// ObjectPrinter is an interface that knows how to print objects.
type ObjectPrinter interface {
	// PrintObj prints the provided object to the provided [io.Writer], according to the printer's
	// configuration.
	//
	// See [printers] for more information.
	PrintObj(obj any, w io.Writer) error
}

// PrintOptions is a struct that holds options for printing objects.
type PrintOptions struct {
	// NoHeaders configures table-based printers to skip printing headers to the output.
	NoHeaders bool
	// Wide configures table-based printers to include additional columns in their output.
	Wide bool
}
