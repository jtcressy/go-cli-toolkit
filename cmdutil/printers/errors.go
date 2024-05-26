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
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/samber/lo"
)

type NoCompatiblePrinterError struct {
	OutputFormat   *string
	AllowedFormats []string
	Options        any
}

func (n NoCompatiblePrinterError) Error() string {
	sort.Strings(n.AllowedFormats)
	return fmt.Sprintf(
		"no compatible printer found for output format %q. Allowed formats: %s",
		lo.FromPtrOr(n.OutputFormat, ""),
		strings.Join(n.AllowedFormats, ", "),
	)
}

// IsNoCompatiblePrinterError returns true if the error is a NoCompatiblePrinterError.
func IsNoCompatiblePrinterError(err error) bool {
	var n NoCompatiblePrinterError
	return errors.As(err, &n)
}
