package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jtpck/installer/config"
	"github.com/jtpck/installer/shell"
	"github.com/jtpck/installer/ui"
	"github.com/jtpck/installer/validator"
	"github.com/jtpck/installer/wrapper"
	"github.com/spf13/cobra"
)

var (
	demoMode bool
	version  = "0.1.0"
)

const endpoint = "https://JTPCK.com/api/v1/telemetry"

var uuidRegex = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

func validateUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}

var rootCmd = &cobra.Command{
	Use:   "jtpck [user_id]",
	Short: "JTPCK Telemetry Setup for Claude Code",
	Long: `Interactive setup tool for configuring OpenTelemetry for Claude Code, Codex, and the Gemini CLI.

Provide user_id as argument for one-command setup:
  jtpck abc123

Or run interactively:
  jtpck`,
	Args: cobra.MaximumNArgs(1),
	Run:  runSetup,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVar(&demoMode, "demo", false, "Demo mode (UI preview without file writes)")
	rootCmd.Version = version
}

func runSetup(cmd *cobra.Command, args []string) {
	tools := []string{"claude", "codex", "gemini"}

	// Check if user ID provided as argument
	var userID string
	if len(args) > 0 {
		userID = strings.ToLower(strings.TrimSpace(args[0]))
		// Validate UUID format for command-line args
		if !validateUUID(userID) {
			fmt.Printf("Error: Invalid user ID format. Must be a valid UUID (e.g., 12345678-1234-1234-1234-123456789abc)\n")
			os.Exit(1)
		}
	}

	// In demo mode, skip validation and config checks
	if !demoMode {
		// Check if already configured
		if config.Exists() && userID == "" {
			fmt.Println("⚠ Configuration already exists.")
			fmt.Print("Reconfigure? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Setup cancelled.")
				return
			}
		}

		// Validate tools
		missing := validator.GetMissingTools(tools)
		if len(missing) > 0 {
			fmt.Printf("⚠ Warning: The following tools are not installed: %v\n", missing)
			fmt.Println("Wrappers will only be created for installed tools.\n")
		}
	}

	// Run animation
	animModel := ui.NewAnimationModel()
	p := tea.NewProgram(animModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running animation: %v\n", err)
		os.Exit(1)
	}

	// Run input screen only if user ID not provided
	if userID == "" {
		var currentValue string
		if !demoMode && config.Exists() {
			cfg, _ := config.Load()
			if cfg != nil {
				currentValue = cfg.UserID
			}
		}

		inputModel := ui.NewInputModel(currentValue)
		p = tea.NewProgram(inputModel)
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running input: %v\n", err)
			os.Exit(1)
		}

		inputResult := finalModel.(ui.InputModel)
		if !inputResult.Done() {
			fmt.Println("Setup cancelled.")
			return
		}

		userID = inputResult.GetUserID()
	}

	var actions []string

	// Generate per-application env vars
	appEnvs := config.AppEnvs(userID, endpoint)

	// Add Claude Code action
	actions = append(actions, "Enabled Claude Code telemetry (wrapper script)")

	// Configure Codex telemetry config file (respects demo mode)
	actions = append(actions, fmt.Sprintf("Enabled Codex telemetry (config at %s)", config.CodexConfigPath()))
	if err := config.EnableCodexTelemetry(userID, endpoint, demoMode, func(msg string) { fmt.Println(msg) }); err != nil {
		fmt.Printf("Error enabling Codex telemetry: %v\n", err)
		os.Exit(1)
	}

	// Configure Gemini telemetry settings/file (respects demo mode)
	actions = append(actions, fmt.Sprintf("Enabled Gemini CLI telemetry (settings at %s)", config.GeminiSettingsPath()))
	if err := config.EnableGeminiTelemetry(userID, endpoint, demoMode, func(msg string) { fmt.Println(msg) }); err != nil {
		fmt.Printf("Error enabling Gemini telemetry: %v\n", err)
		os.Exit(1)
	}

	var installedTools []string

	if !demoMode {
		// Save config
		cfg := config.New(userID, endpoint)
		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}

		// Create wrappers only for installed tools
		installed := validator.GetInstalledTools(tools)
		if err := wrapper.CreateWrappers(appEnvs, installed); err != nil {
			fmt.Printf("Error creating wrappers: %v\n", err)
			os.Exit(1)
		}

		installedTools = wrapper.GetInstalledTools(tools)

		// Auto-install aliases to shell config
		if err := shell.InstallAliases(installedTools); err != nil {
			fmt.Printf("Warning: Could not auto-install aliases: %v\n", err)
			fmt.Println("You'll need to manually add aliases to your shell config.")
		}
	} else {
		// In demo mode, pretend all tools are installed
		installedTools = tools
	}

	// Detect shell config
	shellConfig := shell.DetectShellConfig()
	aliasCommands := shell.GenerateAliasCommands(installedTools)

	// Run success screen (always show auto-installed UI)
	successModel := ui.NewSuccessModel(shellConfig, installedTools, aliasCommands, true, actions)
	p = tea.NewProgram(successModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running success screen: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ JTPCK setup complete!")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  🔄 Run now: \033[1;36msource ~/%s\033[0m\n", shellConfig)
	fmt.Println("  Or restart your terminal")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
