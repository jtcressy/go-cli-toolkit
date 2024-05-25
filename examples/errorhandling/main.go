package main

import (
	"errors"
	"os"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

func main() {
	defer err2.Catch(func(err error) error {
		// Use this to catch and handle errors
		os.Exit(1)
		return nil
	}, func(p any) {
		// Use this to handle panics
		os.Exit(1)
	})

	// Do something that errors
	try.To(func() error { return errors.New("something went wrong") }())
}
