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
	"github.com/jtcressy/go-cli-toolkit/cmdutil/printers"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TableCSVPrinterFlags", Label("unit"), func() {
	var (
		tableCSVPrinterFlags *printers.TableCSVPrinterFlags
	)

	BeforeEach(func() {
		tableCSVPrinterFlags = &printers.TableCSVPrinterFlags{}
	})

	Describe("AddFlags", func() {
		It("should add no-headers flag when NoHeaders is not nil", func() {
			noHeaders := false
			tableCSVPrinterFlags.NoHeaders = &noHeaders
			cmd := &cobra.Command{}
			tableCSVPrinterFlags.AddFlags(cmd)
			flag := cmd.Flag("no-headers")
			Expect(flag).ToNot(BeNil())
			Expect(flag.Value.String()).To(Equal("false"))
		})

		It("should not add no-headers flag when NoHeaders is nil", func() {
			tableCSVPrinterFlags.NoHeaders = nil
			cmd := &cobra.Command{}
			tableCSVPrinterFlags.AddFlags(cmd)
			flag := cmd.Flag("no-headers")
			Expect(flag).To(BeNil())
		})
	})

	Describe("AllowedFormats", func() {
		It("should return csv and table as allowed formats", func() {
			formats := tableCSVPrinterFlags.AllowedFormats()
			Expect(formats).To(ConsistOf("csv", "table"))
		})
	})

	Describe("ToPrinter", func() {
		It("should return CSVPrinter when format is csv", func() {
			printer, err := tableCSVPrinterFlags.ToPrinter("csv")
			Expect(err).ToNot(HaveOccurred())
			Expect(printer).To(BeAssignableToTypeOf(&printers.CSVPrinter{}))
		})

		It("should return TablePrinter when format is table", func() {
			printer, err := tableCSVPrinterFlags.ToPrinter("table")
			Expect(err).ToNot(HaveOccurred())
			Expect(printer).To(BeAssignableToTypeOf(&printers.TablePrinter{}))
		})

		It("should return error when format is not supported", func() {
			printer, err := tableCSVPrinterFlags.ToPrinter("unsupported")
			Expect(err).To(HaveOccurred())
			Expect(printer).To(BeNil())
			Expect(err).To(MatchError(printers.NoCompatiblePrinterError{
				OutputFormat:   lo.ToPtr("unsupported"),
				AllowedFormats: []string{"csv", "table"},
			}))
		})
	})
})
