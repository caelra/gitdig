package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitdig",
	Short: "creates github reports",
}

func Execute() error {
	return rootCmd.Execute()
}
