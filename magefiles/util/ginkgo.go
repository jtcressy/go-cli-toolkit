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

package util

import (
	"fmt"

	"github.com/jtcressy/go-cli-toolkit/magefiles/helpers"
	"github.com/magefile/mage/mg"
)

var (
	ginkgo = helpers.GoRunCmdV("github.com/onsi/ginkgo/v2/ginkgo")
)

type Ginkgo mg.Namespace

type GinkgoOptions struct {
	PackagePath    string
	LabelFilter    string
	CoverageReport string
	JSONReport     string
	JunitReport    string
	CI             bool
	Watch          bool
}

// Bootstrap bootstraps a ginkgo suite at the given path.
func (u Ginkgo) Bootstrap(path string) error {
	fmt.Printf("Bootstrapping ginkgo suite (at %s)...\n", path)
	return helpers.ExecuteInDirectory(path, func(_ ...string) error {
		return ginkgo("bootstrap")
	}, false)
}

func (u Ginkgo) Run(opts GinkgoOptions) error {
	cmd := "run"
	if opts.Watch {
		cmd = "watch"
	}
	args := []string{cmd,
		"-r", "--randomize-all",
		"--trace", "--fail-on-pending", fmt.Sprintf("--label-filter=%s", opts.LabelFilter),
	}
	if !opts.Watch {
		args = append(args,
			"--randomize-suites",
			"--keep-going",
		)
	}
	if opts.CI {
		args = append(
			args,
			"--no-color",
			"--github-output",
		)
	}
	if opts.CoverageReport != "" {
		args = append(
			args,
			"--cover",
			fmt.Sprintf("--coverprofile=%s", opts.CoverageReport),
			fmt.Sprintf("--coverpkg=%s", opts.PackagePath),
		)
	}
	if opts.JSONReport != "" {
		args = append(args, fmt.Sprintf("--json-report=%s", opts.JSONReport))
	}
	if opts.JunitReport != "" {
		args = append(args, fmt.Sprintf("--junit-report=%s", opts.JunitReport))
	}
	args = append(args, opts.PackagePath)
	return ginkgo(args...)
}
