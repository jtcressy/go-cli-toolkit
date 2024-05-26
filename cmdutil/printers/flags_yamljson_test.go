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

var _ = Describe("YamlJSONPrinterFlags", Label("unit"), func() {
	var (
		cmd *cobra.Command
		y   *printers.YamlJSONPrinterFlags
	)

	BeforeEach(func() {
		cmd = &cobra.Command{}
		y = &printers.YamlJSONPrinterFlags{}
	})

	Describe("AddFlags", func() {
		Context("when JSONIndent is not nil", func() {
			It("adds a json-indent flag to the command", func() {
				indent := true
				y.JSONIndent = &indent
				y.AddFlags(cmd)
				Expect(cmd.Flags().Lookup("json-indent")).NotTo(BeNil())
			})
		})

		Context("when JSONIndent is nil", func() {
			It("does not add a json-indent flag to the command", func() {
				y.JSONIndent = nil
				y.AddFlags(cmd)
				Expect(cmd.Flags().Lookup("json-indent")).To(BeNil())
			})
		})
	})

	Describe("AllowedFormats", func() {
		It("returns json and yaml as allowed formats", func() {
			Expect(y.AllowedFormats()).To(Equal([]string{"json", "yaml"}))
		})
	})

	Describe("ToPrinter", func() {
		Context("when format is json", func() {
			It("returns a JSONPrinter", func() {
				printer, err := y.ToPrinter("json")
				Expect(err).ToNot(HaveOccurred())
				Expect(printer).To(BeAssignableToTypeOf(&printers.JSONPrinter{}))
			})
		})

		Context("when format is yaml", func() {
			It("returns a YAMLPrinter", func() {
				printer, err := y.ToPrinter("yaml")
				Expect(err).ToNot(HaveOccurred())
				Expect(printer).To(BeAssignableToTypeOf(&printers.YamlPrinter{}))
			})
		})

		Context("when format is not json or yaml", func() {
			It("returns an error", func() {
				_, err := y.ToPrinter("xml")
				Expect(err).To(MatchError(printers.NoCompatiblePrinterError{
					OutputFormat:   lo.ToPtr("xml"),
					AllowedFormats: []string{"json", "yaml"},
				}))
			})
		})
	})
})
