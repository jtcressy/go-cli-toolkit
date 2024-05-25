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

package flags_test

import (
	"github.com/jtcressy/go-cli-toolkit/flags"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlagProxy", func() {
	var (
		flagProxy *flags.FlagProxy
	)

	BeforeEach(func() {
		flagProxy = &flags.FlagProxy{
			StringFn: func() string {
				return "value"
			},
			SetFn: func(_ string) error {
				return nil
			},
			TypeFn: func() string {
				return "string"
			},
		}
	})

	Describe("String", func() {
		It("should return the string representation of the value", func() {
			Expect(flagProxy.String()).To(Equal("value"))
		})
	})

	Describe("Set", func() {
		It("should set the value from the string representation", func() {
			err := flagProxy.Set("new value")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Type", func() {
		It("should return the type of the value", func() {
			Expect(flagProxy.Type()).To(Equal("string"))
		})
	})
})
