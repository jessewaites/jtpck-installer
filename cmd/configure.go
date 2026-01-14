package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jtpck/installer/config"
	"github.com/jtpck/installer/shell"
	"github.com/jtpck/installer/ui"
	"github.com/jtpck/installer/validator"
	"github.com/jtpck/installer/wrapper"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Reconfigure JTPCK telemetry settings",
	Long:  `Update your user ID and regenerate wrapper scripts.`,
	Run:   runConfigure,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func runConfigure(cmd *cobra.Command, args []string) {
	tools := []string{"claude", "codex", "gemini"}

	// Load existing config
	var currentValue string
	if config.Exists() {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
		currentValue = cfg.UserID
	}

	// Run input screen (skip animation for reconfigure)
	inputModel := ui.NewInputModel(currentValue)
	p := tea.NewProgram(inputModel)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running input: %v\n", err)
		os.Exit(1)
	}

	inputResult := finalModel.(ui.InputModel)
	if !inputResult.Done() {
		fmt.Println("Reconfiguration cancelled.")
		return
	}

	userID := inputResult.GetUserID()

	var actions []string

	// Generate per-application env vars
	appEnvs := config.AppEnvs(userID, endpoint)

	// Add Claude Code action
	actions = append(actions, "Enabled Claude Code telemetry (wrapper script)")

	// Configure Codex telemetry config file
	actions = append(actions, fmt.Sprintf("Enabled Codex telemetry (config at %s)", config.CodexConfigPath()))
	if err := config.EnableCodexTelemetry(userID, endpoint, false, func(msg string) { fmt.Println(msg) }); err != nil {
		fmt.Printf("Error enabling Codex telemetry: %v\n", err)
		os.Exit(1)
	}

	// Configure Gemini telemetry settings/file
	actions = append(actions, fmt.Sprintf("Enabled Gemini CLI telemetry (settings at %s)", config.GeminiSettingsPath()))
	if err := config.EnableGeminiTelemetry(userID, endpoint, false, func(msg string) { fmt.Println(msg) }); err != nil {
		fmt.Printf("Error enabling Gemini telemetry: %v\n", err)
		os.Exit(1)
	}

	// Save updated config
	cfg := config.New(userID, endpoint)
	if err := cfg.Save(); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	// Recreate wrappers
	installed := validator.GetInstalledTools(tools)
	if err := wrapper.CreateWrappers(appEnvs, installed); err != nil {
		fmt.Printf("Error creating wrappers: %v\n", err)
		os.Exit(1)
	}

	installedTools := wrapper.GetInstalledTools(tools)

	// Auto-install aliases to shell config
	if err := shell.InstallAliases(installedTools); err != nil {
		fmt.Printf("Warning: Could not auto-install aliases: %v\n", err)
		fmt.Println("You'll need to manually add aliases to your shell config.")
	}

	// Detect shell config
	shellConfig := shell.DetectShellConfig()
	aliasCommands := shell.GenerateAliasCommands(installedTools)

	// Run success screen
	successModel := ui.NewSuccessModel(shellConfig, installedTools, aliasCommands, true, actions)
	p = tea.NewProgram(successModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running success screen: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ Reconfiguration complete!")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  🔄 Run now: \033[1;36msource ~/%s\033[0m\n", shellConfig)
	fmt.Println("  Or restart your terminal")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
