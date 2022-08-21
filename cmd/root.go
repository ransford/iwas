package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "iwas",
	Short: "Versioned policy browser for AWS IAM policies",
	Long:  "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setLogLevel(c *cobra.Command, args []string) {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debugging output")
}
