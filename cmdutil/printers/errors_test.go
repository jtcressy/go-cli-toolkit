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
	"errors"

	"github.com/jtcressy/go-cli-toolkit/cmdutil/printers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error handling in printers", Label("unit"), func() {
	Describe("NoCompatiblePrinterError", func() {
		It("should return the correct error message", func() {
			err := printers.NoCompatiblePrinterError{
				OutputFormat:   nil,
				AllowedFormats: []string{"json", "yaml"},
				Options:        nil,
			}
			Expect(
				err.Error(),
			).To(Equal("no compatible printer found for output format \"\". Allowed formats: json, yaml"))
		})

		It("should return the correct error message with output format", func() {
			format := "xml"
			err := printers.NoCompatiblePrinterError{
				OutputFormat:   &format,
				AllowedFormats: []string{"json", "yaml"},
				Options:        nil,
			}
			Expect(
				err.Error(),
			).To(Equal("no compatible printer found for output format \"xml\". Allowed formats: json, yaml"))
		})
	})

	Describe("IsNoCompatiblePrinterError", func() {
		It("should return true for NoCompatiblePrinterError", func() {
			err := printers.NoCompatiblePrinterError{
				OutputFormat:   nil,
				AllowedFormats: []string{"json", "yaml"},
				Options:        nil,
			}
			Expect(printers.IsNoCompatiblePrinterError(err)).To(BeTrue())
		})

		It("should return false for other errors", func() {
			err := errors.New("some other error")
			Expect(printers.IsNoCompatiblePrinterError(err)).To(BeFalse())
		})
	})
})
