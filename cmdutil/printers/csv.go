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
	"encoding/csv"
	"io"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

var _ ObjectPrinter = (*CSVPrinter)(nil)

type CSVPrinter struct {
	PrintOptions
	TableReflectorFunc
}

func (p *CSVPrinter) PrintObj(a any, w io.Writer) (err error) {
	defer err2.Handle(&err, nil)
	enc := csv.NewWriter(w)
	headers, rows := try.To2(p.TableReflectorFunc(a))
	if !p.NoHeaders {
		try.To(enc.Write(headers))
	}
	try.To(enc.WriteAll(rows))
	enc.Flush()
	return enc.Error()
}

func NewCSVPrinter(options PrintOptions) ObjectPrinter {
	printer := &CSVPrinter{
		PrintOptions:       options,
		TableReflectorFunc: DefaultTableReflectorFunc,
	}
	return printer
}
