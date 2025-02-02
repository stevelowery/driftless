package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stevelowery/driftless/internal/cmd/squarespace"
	"github.com/stevelowery/driftless/internal/cmd/version"
)

func New() *cobra.Command {

	root := &cobra.Command{
		Use:   "driftless",
		Short: "driftless",
		Long:  "utilities for managing Driftless MTB",
	}

	root.AddCommand(squarespace.New())
	root.AddCommand(version.New())

	return root
}
