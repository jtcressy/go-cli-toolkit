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

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/carolynvs/magex/pkg"
	"github.com/magefile/mage/sh"
)

var (
	// tools.
	golangcilint = "github.com/golangci/golangci-lint/cmd/golangci-lint"
	addlicense   = "github.com/google/addlicense"
	ginkgo       = "github.com/onsi/ginkgo/v2/ginkgo"
	golines      = "github.com/segmentio/golines"
	mockery      = "github.com/vektra/mockery/v2"

	ciTools  = []string{golangcilint, addlicense, ginkgo, golines}
	allTools = append(ciTools, []string{mockery}...)
)

// Setup runs the setup script for the current OS.
func main() {
	var err error

	// Ensure Mage is installed and available on the $PATH.
	if err = pkg.EnsureMage(""); err != nil {
		panic(err)
	}

	var toInstall []string
	if isCi() {
		toInstall = ciTools
	} else {
		toInstall = allTools
	}

	if err = setupGoDeps(toInstall); err != nil {
		panic(err)
	}
}

func isCi() bool {
	return len(os.Args) > 1 && os.Args[1] == "ci"
}

func setupGoDeps(toInstall []string) error {
	for _, tool := range toInstall {
		fmt.Println("Installing", fmt.Sprintf("`%s`", tool))
		if err := sh.RunCmd("go", "install", "-mod=readonly", tool); err() != nil {
			return errors.New("failed to install " + tool + ": " + err().Error())
		}
	}
	fmt.Println("\n==============================================================")
	fmt.Println("All Tools installed successful! Ensure $GOPATH/bin is on your $PATH!")
	fmt.Println("==============================================================")
	return nil
}
