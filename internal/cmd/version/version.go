package version

import "github.com/spf13/cobra"

var (
	Version = "development"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(Version)
		},
	}
}
