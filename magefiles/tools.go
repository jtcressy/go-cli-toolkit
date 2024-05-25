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

var (
	goGenerate = helpers.RunCmdV("go", "generate")
	goModTidy  = helpers.RunCmdV("go", "mod", "tidy")
	goWorkSync = helpers.RunCmdV("go", "work", "sync")
	gitDiff    = helpers.RunCmdV(
		"git",
		"diff",
		"--stat",
		"--exit-code",
		".",
		"':(exclude)*.mod' ':(exclude)*.sum'",
	)
)

// Runs `go generate` on the entire project.
func Generate() error {
	// if err := mockery(); err != nil {
	// 	return err
	// }
	return helpers.ExecuteForAllModules(
		helpers.ListGoModulesFromGoWork("go.work"),
		func(...string) error { return goGenerate("./...") },
		false,
	)
}

// Runs `go generate` on the entire project and verifies that no files were
// changed.
func GenerateCheck() error {
	if err := helpers.ExecuteForAllModules(
		helpers.ListGoModulesFromGoWork("go.work"), func(...string) error { return goGenerate("./...") }, false,
	); err != nil {
		return err
	}
	if err := gitDiff(); err != nil {
		return fmt.Errorf("generated files are out of date: %w", err)
	}
	return nil
}

// Runs 'go tidy' on the entire project.
func Tidy() error {
	return helpers.ExecuteForAllModules(
		helpers.ListGoModulesFromGoWork("go.work"),
		goModTidy,
		false,
	)
}

// Runs 'go work sync' on the entire project.
func Sync() error {
	return goWorkSync()
}

// Cleans the project.
func Clean() error {
	// Remove golang build artifacts.
	if err := sh.Rm("bin"); err != nil {
		return err
	}

	// Remove test cache.
	return sh.RunV("go", "clean", "-testcache")
}
