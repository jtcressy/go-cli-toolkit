package util

import "github.com/spf13/cobra"

func CobraCMDErr2HandlerAnnotation(cmd *cobra.Command) (s string) {
	s = cmd.Name()
	cmd.VisitParents(func(c *cobra.Command) {
		s = c.Name() + " " + s
	})
	return s
}
