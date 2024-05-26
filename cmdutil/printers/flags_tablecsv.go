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
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var _ FlaggablePrinter = (*TableCSVPrinterFlags)(nil)

type TableCSVPrinterFlags struct {
	NoHeaders *bool
}

// AddFlags implements FlaggablePrinter.
func (t *TableCSVPrinterFlags) AddFlags(cmd *cobra.Command) {
	if t.NoHeaders != nil {
		cmd.Flags().BoolVar(
			t.NoHeaders,
			"no-headers",
			lo.FromPtrOr(t.NoHeaders, false),
			"When using the default table output, don't print headers (default print headers).",
		)
	}
}

// AllowedFormats implements FlaggablePrinter.
func (t *TableCSVPrinterFlags) AllowedFormats() []string {
	return []string{"csv", "table"}
}

// ToPrinter implements FlaggablePrinter.
func (t *TableCSVPrinterFlags) ToPrinter(format string) (ObjectPrinter, error) {
	switch format {
	case "csv":
		return NewCSVPrinter(PrintOptions{
			NoHeaders: lo.FromPtrOr(t.NoHeaders, false),
		}), nil
	case "table":
		return NewTablePrinter(PrintOptions{
			NoHeaders: lo.FromPtrOr(t.NoHeaders, false),
		}), nil
	default:
		return nil, NoCompatiblePrinterError{
			OutputFormat:   lo.ToPtr(format),
			AllowedFormats: t.AllowedFormats(),
		}
	}
}
