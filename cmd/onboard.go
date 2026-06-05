// cmd/onboard.go
package cmd

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Configure baltig with your GitLab instance",
	RunE:  runOnboard,
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}

func runOnboard(cmd *cobra.Command, args []string) error {
	fmt.Println("baltig onboard")
	fmt.Println()

	gitlabURL := promptLine("GitLab URL", "https://gitlab.example.com")
	gitlabURL = strings.TrimRight(gitlabURL, "/")

	token := promptMasked("Personal Access Token (scope: api)")

	fmt.Print("\nConnecting... ")
	client, err := gitlab.New(gitlabURL, token)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	username, err := client.Ping()
	if err != nil {
		fmt.Println("✗ Failed")
		return fmt.Errorf("connect to %s: %w", gitlabURL, err)
	}
	fmt.Printf("✓ Connected as @%s\n\n", username)

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	cfg.Global.GitLabURL = gitlabURL
	cfg.Global.Token = token
	if cfg.Global.DefaultRef == "" {
		cfg.Global.DefaultRef = "main"
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	cfgPath, err := config.ConfigPath()
	if err != nil {
		cfgPath = "~/.baltig/config.yaml"
	}
	fmt.Printf("Config saved to %s\n", cfgPath)
	fmt.Println("Run 'baltig' to start.")
	return nil
}

func promptLine(label, placeholder string) string {
	fmt.Printf("%s [%s]: ", label, placeholder)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return placeholder
	}
	return input
}

func promptMasked(label string) string {
	fmt.Printf("%s: ", label)
	bytePass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		var input string
		fmt.Scanln(&input)
		return input
	}
	return string(bytePass)
}

