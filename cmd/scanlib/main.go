package main

import (
	"os"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/eval"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scanlib <scanspec> [input]",
		Short: "Evaluate a Scanspec file against input",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			specsrc, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			source, err := ast.ParseString(args[0], string(specsrc))
			if err != nil {
				return err
			}

			input := os.Stdin
			if len(args) >= 2 {
				input, err = os.Open(args[1])
				if err != nil {
					return err
				}
				defer input.Close()
			}

			_, err = eval.Evaluate(source, input)
			if err != nil {
				return err
			}

			return nil
		},
		SilenceUsage: true,
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
