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
	"fmt"

	"github.com/jtcressy/go-cli-toolkit/magefiles/helpers"
	"github.com/magefile/mage/sh"
)

// Run all linters.
func Lint() error {
	cmds := []func() error{GoFmtCheck, GolangCiLint, LicenseCheck}
	for _, cmd := range cmds {
		if err := cmd(); err != nil {
			return err
		}
	}
	return nil
}

// Run all formatters.
func Format() error {
	cmds := []func() error{GoFmt, Golines, License, GolangCiLintFix}
	for _, cmd := range cmds {
		if err := cmd(); err != nil {
			return err
		}
	}
	return nil
}

// Run `go fmt`.
func GoFmt() error {
	return helpers.ExecuteForAllModules(
		helpers.ListGoModulesFromGoWork("go.work"),
		func(_ ...string) error {
			return sh.RunV("go", "fmt", "./...")
		},
		false,
	)
}

// Run `go fmt -n`.
func GoFmtCheck() error {
	o, err := sh.Output("gofmt", "-l", ".")
	if err != nil {
		return err
	}
	if len(o) > 0 {
		return fmt.Errorf("unformatted files: %s", o)
	}
	return nil
}

// Run `golangci-lint`.
func GolangCiLint() error {
	for _, dir := range helpers.ListGoModulesFromGoWork("go.work") {
		if err := helpers.GoRunV("github.com/golangci/golangci-lint/cmd/golangci-lint",
			"run", "--timeout=10m", "--concurrency", "4", "--config=.golangci.yaml", "-v", "./"+dir+"/"+"...",
		); err != nil {
			return err
		}
	}
	return nil
}

// Run `golangci-lint` with --fix.
func GolangCiLintFix() error {
	for _, dir := range helpers.ListGoModulesFromGoWork("go.work") {
		if err := helpers.GoRunV("github.com/golangci/golangci-lint/cmd/golangci-lint",
			"run", "--timeout=10m", "--concurrency", "4", "--config=.golangci.yaml", "-v", "--fix", "./"+dir+"/"+"...",
		); err != nil {
			return err
		}
	}
	return nil
}

// Run `golines`.
func Golines() error {
	return helpers.GoRunV("github.com/segmentio/golines",
		"--reformat-tags", "--shorten-comments", "--write-output", "--max-len=99", "-l", "./.",
	)
}

// Run `addlicense`.
func License() error {
	return helpers.ExecuteForAllModules(
		helpers.ListGoModulesFromGoWork("go.work"),
		func(_ ...string) error {
			return helpers.GoRunV("github.com/google/addlicense",
				"-v", "-f", "./LICENSE.header", "./.",
			)
		},
		true,
	)
}

// Run `addlicense` with -check.
func LicenseCheck() error {
	return helpers.ExecuteForAllModules(
		helpers.ListGoModulesFromGoWork("go.work"),
		func(_ ...string) error {
			return helpers.GoRunV("github.com/google/addlicense",
				"-check", "-v", "-f", "./LICENSE.header", "-ignore", "**/*.yml", "./.",
			)
		},
		true,
	)
}
