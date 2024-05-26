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

var _ FlaggablePrinter = (*YamlJSONPrinterFlags)(nil)

type YamlJSONPrinterFlags struct {
	JSONIndent *bool
}

// AddFlags implements FlaggablePrinter.
func (y *YamlJSONPrinterFlags) AddFlags(cmd *cobra.Command) {
	if y.JSONIndent != nil {
		cmd.Flags().
			BoolVar(
				y.JSONIndent,
				"json-indent",
				lo.FromPtrOr(y.JSONIndent, false),
				"When using the \"json\" output format, indent it for better readability (default no indent).",
			)
	}
}

// AllowedFormats implements FlaggablePrinter.
func (y *YamlJSONPrinterFlags) AllowedFormats() []string {
	return []string{"json", "yaml"}
}

// ToPrinter implements FlaggablePrinter.
func (y *YamlJSONPrinterFlags) ToPrinter(format string) (ObjectPrinter, error) {
	switch format {
	case "json":
		return NewJSONPrinter(lo.FromPtrOr(y.JSONIndent, false)), nil
	case "yaml":
		return NewYAMLPrinter(), nil
	default:
		return nil, NoCompatiblePrinterError{
			OutputFormat:   lo.ToPtr(format),
			AllowedFormats: y.AllowedFormats(),
		}
	}
}
