// cmd/doctor.go
package cmd

import (
	"fmt"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check baltig configuration and GitLab connectivity",
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("baltig doctor")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		printCheck(false, "config", err.Error())
		return nil
	}

	cfgPath, err := config.ConfigPath()
	if err != nil {
		printCheck(false, "config path", err.Error())
		return nil
	}
	printCheck(true, "config", cfgPath)

	if err := config.Validate(cfg); err != nil {
		printCheck(false, "config valid", err.Error())
		return nil
	}
	printCheck(true, "config valid", "")

	client, err := gitlab.New(cfg.Global.GitLabURL, cfg.Global.Token)
	if err != nil {
		printCheck(false, "gitlab client", err.Error())
		return nil
	}
	username, err := client.Ping()
	if err != nil {
		printCheck(false, "gitlab connection", err.Error())
		return nil
	}
	printCheck(true, "gitlab connection", fmt.Sprintf("@%s at %s", username, cfg.Global.GitLabURL))

	fmt.Println()
	fmt.Println("All checks passed.")
	return nil
}

func printCheck(ok bool, label, detail string) {
	icon := "✓"
	if !ok {
		icon = "✗"
	}
	if detail != "" {
		fmt.Printf("  %s %s — %s\n", icon, label, detail)
	} else {
		fmt.Printf("  %s %s\n", icon, label)
	}
}
