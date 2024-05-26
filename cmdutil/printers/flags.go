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
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// FlaggablePrinter is used to provision a new group of printer implementations that may add custom
// flags in AddFlags and produce a ObjectPrinter via ToPrinter based on the allowed formats
// returned by AllowedFormats.
type FlaggablePrinter interface {
	// AllowedFormats should return a list of format options that are supported by ToPrinte
	AllowedFormats() []string
	// ToPrinter should return a ObjectPrinter for a given output format string.
	// If the output format string is not valid, or cannot produce a ObjectPrinter,
	// a NoCompatiblePrinterError must be returned.
	ToPrinter(format string) (ObjectPrinter, error)
	// AddFlags should call the given cobra.Command's Flags() method to add new flags to its
	// pflag.FlagSet that may be specific to the available printer implementations.
	AddFlags(cmd *cobra.Command)
}

// PrintFlags composes common printer types
// used across all commands, and provides a method
// of retreiving a known printer based on the flag values provided.
type PrintFlags struct {
	RegisteredPrintFlaggers []FlaggablePrinter
	OutputFormat            *string

	// OutputFlagSpecified indicates whether the user specifically requested a certain kind of
	// output. using this function allows a sophisticated caller to change the flag binding logic
	// if they so desire.
	OutputFlagSpecified func() bool
}

// AllowedFormats returns string slice of allowed printing formats.
func (f *PrintFlags) AllowedFormats() []string {
	return lo.Reduce(
		f.RegisteredPrintFlaggers,
		func(agg []string, rp FlaggablePrinter, _ int) []string {
			return append(agg, rp.AllowedFormats()...)
		},
		make([]string, 0),
	)
}

func (f *PrintFlags) ToPrinter() (ObjectPrinter, error) {
	for _, fp := range f.RegisteredPrintFlaggers {
		if p, err := fp.ToPrinter(lo.FromPtrOr(f.OutputFormat, "")); !IsNoCompatiblePrinterError(
			err,
		) {
			return p, err
		}
	}
	return nil, NoCompatiblePrinterError{
		OutputFormat:   lo.ToPtr(lo.FromPtrOr(f.OutputFormat, "")),
		AllowedFormats: f.AllowedFormats(),
	}
}

// AddFlags takes a *cobra.Command by reference and binds
// flags related to printing to the command.
func (f *PrintFlags) AddFlags(cmd *cobra.Command) {
	lo.ForEach(f.RegisteredPrintFlaggers, func(rp FlaggablePrinter, _ int) {
		rp.AddFlags(cmd)
	})
	if f.OutputFormat != nil {
		cmd.Flags().StringVarP(
			f.OutputFormat,
			"output",
			"o",
			lo.FromPtrOr(f.OutputFormat, ""),
			fmt.Sprintf("Output format. One of: (%s).", strings.Join(f.AllowedFormats(), ", ")),
		)
		if f.OutputFlagSpecified == nil {
			f.OutputFlagSpecified = func() bool {
				return cmd.Flag("output").Changed
			}
		}
	}
}

// WithDefaultOutput sets a default output format if one is not provided through a flag value.
func (f *PrintFlags) WithDefaultOutput(format string) *PrintFlags {
	f.OutputFormat = &format
	return f
}

func NewPrintFlags() *PrintFlags {
	return &PrintFlags{
		OutputFormat: lo.ToPtr(""),
		RegisteredPrintFlaggers: []FlaggablePrinter{
			&YamlJSONPrinterFlags{
				JSONIndent: lo.ToPtr(false),
			},
		},
	}
}
