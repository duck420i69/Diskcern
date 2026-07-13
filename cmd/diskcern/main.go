package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/diskcern/diskcern/internal/ui"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "diskcern",
		Short: "Diskcern is a fast disk analysis and snapshot tool",
		Run: func(cmd *cobra.Command, args []string) {
			// If no specific CLI command is given, launch the Gio UI
			ui.RunApp()
		},
	}

	var scanCmd = &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a directory and create a snapshot",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			fmt.Printf("Starting scan for %s...\n", path)
			// TODO: initialize DB, run parallel scanner
		},
	}

	var diffCmd = &cobra.Command{
		Use:   "diff [snapshot1] [snapshot2]",
		Short: "Compare two snapshots",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Comparing snapshots %s and %s...\n", args[0], args[1])
			// TODO: diff snapshots
		},
	}

	rootCmd.AddCommand(scanCmd, diffCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
