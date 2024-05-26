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
	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PrintFlags", Label("unit"), func() {
	var (
		printFlags *printers.PrintFlags
		cmd        *cobra.Command
	)

	BeforeEach(func() {
		printFlags = printers.NewPrintFlags()
		cmd = &cobra.Command{}
	})

	Context("when adding flags", func() {
		It("should add output flag to command", func() {
			printFlags.AddFlags(cmd)
			Expect(cmd.Flag("output")).NotTo(BeNil())
		})

		It("should set OutputFlagSpecified function", func() {
			printFlags.AddFlags(cmd)
			Expect(printFlags.OutputFlagSpecified).NotTo(BeNil())
		})
	})

	Context("when setting default output", func() {
		It("should set default output format", func() {
			printFlags.WithDefaultOutput("json")
			Expect(*printFlags.OutputFormat).To(Equal("json"))
		})
	})

	Context("when getting printer", func() {
		It("should return a printer if format is supported", func() {
			printFlags.WithDefaultOutput("json")
			printer, err := printFlags.ToPrinter()
			Expect(err).NotTo(HaveOccurred())
			Expect(printer).NotTo(BeNil())
		})

		It("should return an error if format is not supported", func() {
			printFlags.WithDefaultOutput("unsupported")
			_, err := printFlags.ToPrinter()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when getting allowed formats", func() {
		It("should return a list of allowed formats", func() {
			formats := printFlags.AllowedFormats()
			Expect(formats).To(ContainElement("json"))
			Expect(formats).To(ContainElement("yaml"))
		})
	})
})
