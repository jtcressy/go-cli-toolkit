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

var _ = Describe("YamlPrinter", Label("unit"), func() {
	var (
		printer printers.ObjectPrinter
		buffer  *bytes.Buffer
	)

	BeforeEach(func() {
		printer = printers.NewYAMLPrinter()
		buffer = &bytes.Buffer{}
	})

	Context("when printing an object", func() {
		It("should print a simple object correctly", func() {
			err := printer.PrintObj(map[string]string{"key": "value"}, buffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(buffer.String()).To(Equal("key: value\n"))
		})

		It("should print a complex object correctly", func() {
			err := printer.PrintObj(
				map[string]interface{}{"key": map[string]string{"nestedKey": "nestedValue"}},
				buffer,
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(buffer.String()).To(Equal("key:\n    nestedKey: nestedValue\n"))
		})

		It("should return an error when the object cannot be marshalled", func() {
			defer func() {
				if r := recover(); r != nil {
					ExpectWithOffset(1, r).To(ContainSubstring("cannot marshal type: chan bool"))
				} else {
					Fail("expected panic")
				}
			}()
			err := printer.PrintObj(make(chan bool), buffer)
			Expect(err).To(HaveOccurred())
		})
	})
})
