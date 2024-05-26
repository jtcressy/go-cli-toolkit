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
	"strings"

	"github.com/charmbracelet/x/exp/teatest"
	"github.com/jtcressy/go-cli-toolkit/cmdutil/printers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ginkgoTBWrapper struct {
	*GinkgoTBWrapper
}

func (tb *ginkgoTBWrapper) Name() string {
	return strings.ReplaceAll(tb.GinkgoTBWrapper.Name(), " ", "_")
}

var _ = Describe("TablePrinter", Label("unit"), func() {
	var (
		printer *printers.TablePrinter
		buffer  *bytes.Buffer
	)

	BeforeEach(func() {
		buffer = new(bytes.Buffer)
		if p, ok := printers.NewTablePrinter(printers.PrintOptions{}).(*printers.TablePrinter); !ok {
			Fail("failed to create TablePrinter")
		} else {
			printer = p
		}
	})

	Describe("PrintObj", func() {
		Context("when given a valid object", func() {
			It("should print the object as a table", func() {
				obj := []string{"Hello", "World"}
				err := printer.PrintObj(obj, buffer)
				Expect(err).NotTo(HaveOccurred())
				GinkgoT()
				teatest.RequireEqualOutput(&ginkgoTBWrapper{GinkgoTB()}, buffer.Bytes())
			})
		})

		Context("when given an invalid object", func() {
			It("should not print anything", func() {
				obj := make(chan int)
				err := printer.PrintObj(obj, buffer)
				Expect(err).NotTo(HaveOccurred())
				Expect(buffer.String()).To(BeEmpty())
			})
		})

		Context("when given a nil object", func() {
			It("should not print anything", func() {
				err := printer.PrintObj(nil, buffer)
				Expect(err).NotTo(HaveOccurred())
				Expect(buffer.String()).To(BeEmpty())
			})
		})
	})
})
