package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error { return nil },
}

func init() {
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
