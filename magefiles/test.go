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

	"github.com/jtcressy/go-cli-toolkit/magefiles/util"
	"github.com/magefile/mage/mg"
)

type Test mg.Namespace

func (t Test) Unit(path string, cover bool) error {
	fmt.Println("Running Unit tests...")
	var coverageReport string
	if cover {
		coverageReport = "coverage.out"
	}
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		PackagePath:    path,
		LabelFilter:    "!acceptance,!integration,!network,!performance",
		CoverageReport: coverageReport,
	})
}

func (t Test) Watch(path string) error {
	fmt.Println("Watching Unit tests...")
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		PackagePath: path,
		LabelFilter: "!acceptance,!integration,!network,!performance",
		Watch:       true,
	})
}

func (t Test) Acceptance(path string) error {
	fmt.Println("Running Acceptance tests...")
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		PackagePath: path,
		LabelFilter: "acceptance",
	})
}

func (t Test) Integration(path string) error {
	fmt.Println("Running Integration tests...")
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		PackagePath: path,
		LabelFilter: "integration",
	})
}

func (t Test) Performance(path string) error {
	fmt.Println("Running Performance tests...")
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		PackagePath: path,
		LabelFilter: "performance",
	})
}

func (t Test) Ci(coverageReport, jsonReport, junitReport string) error {
	fmt.Println("Running CI tests...")
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		PackagePath:    "./...",
		LabelFilter:    "!network",
		CoverageReport: coverageReport,
		JSONReport:     jsonReport,
		JunitReport:    junitReport,
		CI:             true,
	})
}
