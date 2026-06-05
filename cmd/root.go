// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/benly/baltig/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
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
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if err := config.Validate(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "baltig: "+err.Error())
		fmt.Fprintln(os.Stderr, "Run 'baltig onboard' to configure.")
		os.Exit(1)
	}

	client, err := gitlab.New(cfg.Global.GitLabURL, cfg.Global.Token)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	app := tui.NewApp(cfg, client)
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
