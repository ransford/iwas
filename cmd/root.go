package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "iwas",
	Short: "Versioned policy browser for AWS IAM policies",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hi")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
