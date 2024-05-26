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

package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Flaggable interface {
	SetupFlags(*pflag.FlagSet) error
}

var _ pflag.Value = (*FlagProxy)(nil)

// FlagProxy can be used to set up custom flags that may need to use
// getters/setters to control a particular value, instead of using pointer references.
//
// StringFn, SetFn, and TypeFn should each be initialized with an anonymous function
// to satisfy the pflag.Value interface.
//
// Call your value's setter in SetFn and your value's getter in StringFn while also converting
// your value to/from a string representation.
//
// # Example
//
//	var f *pflag.FlagSet
//	var myobject *MyObject
//	myflag := f.VarPF(&FlagProxy{
//		StringFn: func() string {
//			return myobject.GetMyValue()
//		},
//		SetFn: func(s string) error {
//			return myobject.SetMyValue(s)
//		},
//		TypeFn: func() string {
//			return reflect.TypeOf(myobject.GetMyValue()).Name()
//		},
//	}, "my-flag", "f", "My Flag Does Something")
type FlagProxy struct {
	// StringFn should return the string representation of the value.
	// Implements pflag.Value interface
	StringFn func() string
	// SetFn should set the value from the string representation, converting if necessary.
	// Implements pflag.Value interface
	SetFn func(string) error
	// TypeFn should return the type of the value.
	// Implements pflag.Value interface
	TypeFn func() string
}

// String implements pflag.Value interface.
func (f *FlagProxy) String() string {
	if f.StringFn == nil {
		panic("StringFn not set on FlagProxy")
	}
	return f.StringFn()
}

// Set implements pflag.Value interface.
func (f *FlagProxy) Set(s string) error {
	if f.SetFn == nil {
		return fmt.Errorf("SetFn not set on FlagProxy")
	}
	return f.SetFn(s)
}

// Type implements pflag.Value interface.
func (f *FlagProxy) Type() string {
	if f.TypeFn == nil {
		panic("TypeFn not set on FlagProxy")
	}
	return f.TypeFn()
}
