package squarespace

import (
	"github.com/spf13/cobra"
	"github.com/stevelowery/driftless/internal/processor/squarespace"
)

func New() *cobra.Command {
	ss := &cobra.Command{
		Use:   "squarespace",
		Short: "Squarespace commands",
	}

	opts := &squarespace.Options{}
	analyze := &cobra.Command{
		Use: "analyze",
		RunE: func(cmd *cobra.Command, args []string) error {
			return squarespace.New().Process(cmd.Context(), opts)
		},
	}

	analyze.Flags().StringVarP(&opts.OrdersFile, "orders-file", "o", "", "path to orders file")
	analyze.Flags().StringVarP(&opts.DontationsFile, "donations-file", "d", "", "path to donations file")
	analyze.Flags().StringVarP(&opts.RaffleFile, "raffle-file", "r", "", "path to raffle file")
	analyze.Flags().StringVarP(&opts.OutputFile, "output-file", "O", "", "path to output file")
	analyze.Flags().StringVarP(&opts.PayoutsFile, "payouts-file", "p", "", "path to payouts file")

	ss.AddCommand(analyze)
	return ss
}
