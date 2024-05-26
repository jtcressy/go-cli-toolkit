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
	"bytes"

	"github.com/jtcressy/go-cli-toolkit/cmdutil/printers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CSVPrinter", Label("unit"), func() {
	var (
		printer printers.ObjectPrinter
		buffer  *bytes.Buffer
	)

	BeforeEach(func() {
		printer = printers.NewCSVPrinter(printers.PrintOptions{})
		buffer = &bytes.Buffer{}
	})

	Context("when printing an object", func() {
		It("should print a single struct correctly", func() {
			err := printer.PrintObj(struct {
				Key   string
				Value string
			}{
				Key:   "key",
				Value: "value",
			}, buffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(buffer.String()).To(Equal("KEY,VALUE\nkey,value\n"))
		})

		It("should print a struct slice correctly", func() {
			err := printer.PrintObj([]struct {
				Key   string
				Value string
			}{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			}, buffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(buffer.String()).To(Equal("KEY,VALUE\nkey1,value1\nkey2,value2\n"))
		})

		It("should not print an invalid object", func() {
			err := printer.PrintObj(func() {}, buffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(buffer.String()).To(Equal("\n"))
		})
	})
})
