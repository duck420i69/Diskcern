package main

import (
	"fmt"
	"os"
	"strings"

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

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate boilerplate code",
	}

	var generateProviderCmd = &cobra.Command{
		Use:   "provider [name]",
		Short: "Generate a new provider",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if len(name) == 0 {
				return
			}
			structName := strings.ToUpper(name[:1]) + name[1:]
			filename := fmt.Sprintf("internal/providers/%s.go", name)
			
			if _, err := os.Stat(filename); err == nil {
				fmt.Printf("Error: Provider file %s already exists\n", filename)
				os.Exit(1)
			}

			content := fmt.Sprintf(`package providers

type %sProvider struct{}

func (p *%sProvider) ID() string {
	return "%s"
}

func (p *%sProvider) Name() string {
	return "%s Provider"
}

func (p *%sProvider) Detect(path string, isDir bool) (bool, ScanDirective) {
	// TODO: Implement detection logic
	return false, ContinueTraversal
}

func (p *%sProvider) Analyze(path string) (AnalysisResult, error) {
	// TODO: Implement analysis logic
	return AnalysisResult{}, nil
}

func (p *%sProvider) GetCleanupActions(path string) []Action {
	// TODO: Implement cleanup actions
	return nil
}
`, structName, structName, name, structName, structName, structName, structName, structName)

			err := os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				fmt.Printf("Failed to write provider: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully generated provider at %s\n", filename)
		},
	}

	generateCmd.AddCommand(generateProviderCmd)

	rootCmd.AddCommand(scanCmd, diffCmd, generateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
