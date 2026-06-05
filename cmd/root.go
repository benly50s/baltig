// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "baltig",
	Short: "GitLab pipeline TUI",
	Long:  "Manage GitLab pipelines from your terminal.",
	RunE:  runRoot,
}

func init() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	// TUI 실행 — Task 16에서 구현
	fmt.Println("baltig: TUI not yet implemented")
	return nil
}
