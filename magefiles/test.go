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
	"strconv"

	"github.com/jtcressy/go-cli-toolkit/magefiles/helpers"
	"github.com/jtcressy/go-cli-toolkit/magefiles/util"
	"github.com/lainio/err2/try"
	"github.com/magefile/mage/mg"
)

// Test mg.Namespace Environment Variables:
//
//   - TEST_PACKAGE_PATH: The package path to run tests on. Defaults to "./...".
//   - TEST_COVERAGE_REPORT: The path to write the coverage report to.
//   - TEST_JSON_REPORT: The path to write the JSON report to.
//   - TEST_JUNIT_REPORT: The path to write the JUnit report to.
//   - TEST_WATCH: Whether to watch for changes and re-run tests.
//   - TEST_LABELS: The labels to filter tests by. Defaults to "unit".
type Test mg.Namespace

// Unit runs ginkgo tests with unit label.
func (t Test) Unit() error {
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		LabelFilter:    "unit",
		PackagePath:    helpers.GetEnv("TEST_PACKAGE_PATH", "./..."),
		CoverageReport: helpers.GetEnv("TEST_COVERAGE_REPORT", ""),
		JSONReport:     helpers.GetEnv("TEST_JSON_REPORT", ""),
		JunitReport:    helpers.GetEnv("TEST_JUNIT_REPORT", ""),
		Watch:          try.To1(strconv.ParseBool(helpers.GetEnv("TEST_WATCH", "false"))),
	})
}

// Acceptance runs ginkgo tests with acceptance label.
func (t Test) Acceptance() error {
	fmt.Println("Running Acceptance tests...")
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		LabelFilter:    "acceptance",
		PackagePath:    helpers.GetEnv("TEST_PACKAGE_PATH", "./..."),
		CoverageReport: helpers.GetEnv("TEST_COVERAGE_REPORT", ""),
		JSONReport:     helpers.GetEnv("TEST_JSON_REPORT", ""),
		JunitReport:    helpers.GetEnv("TEST_JUNIT_REPORT", ""),
		Watch:          try.To1(strconv.ParseBool(helpers.GetEnv("TEST_WATCH", "false"))),
	})
}

// Integration runs ginkgo tests with integration label.
func (t Test) Integration() error {
	fmt.Println("Running Integration tests...")
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		LabelFilter:    "integration",
		PackagePath:    helpers.GetEnv("TEST_PACKAGE_PATH", "./..."),
		CoverageReport: helpers.GetEnv("TEST_COVERAGE_REPORT", ""),
		JSONReport:     helpers.GetEnv("TEST_JSON_REPORT", ""),
		JunitReport:    helpers.GetEnv("TEST_JUNIT_REPORT", ""),
		Watch:          try.To1(strconv.ParseBool(helpers.GetEnv("TEST_WATCH", "false"))),
	})
}

// Performance runs ginkgo tests with performance label.
func (t Test) Performance() error {
	fmt.Println("Running Performance tests...")
	return util.Ginkgo{}.Run(util.GinkgoOptions{
		LabelFilter:    "performance",
		PackagePath:    helpers.GetEnv("TEST_PACKAGE_PATH", "./..."),
		CoverageReport: helpers.GetEnv("TEST_COVERAGE_REPORT", ""),
		JSONReport:     helpers.GetEnv("TEST_JSON_REPORT", ""),
		JunitReport:    helpers.GetEnv("TEST_JUNIT_REPORT", ""),
		Watch:          try.To1(strconv.ParseBool(helpers.GetEnv("TEST_WATCH", "false"))),
	})
}

// Ci runs all ginkgo tests relevant for CI - usage: test:ci "<coverage report path>" "<json report
// path>" "<junit report path>".
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
