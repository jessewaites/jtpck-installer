package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SuccessModel struct {
	shellConfig   string
	tools         []string
	aliasCommands string
	autoInstalled bool
	actions       []string
	done          bool
}

func NewSuccessModel(shellConfig string, tools []string, aliasCommands string, autoInstalled bool, actions []string) SuccessModel {
	return SuccessModel{
		shellConfig:   shellConfig,
		tools:         tools,
		aliasCommands: aliasCommands,
		autoInstalled: autoInstalled,
		actions:       actions,
	}
}

func (m SuccessModel) Init() tea.Cmd {
	return nil
}

func (m SuccessModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m SuccessModel) View() string {
	if m.done {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n")


	if m.autoInstalled {
		// Auto-installed - show simple success message
		sb.WriteString(SuccessStyle.Render("🎉 Success!"))
		sb.WriteString("\n\n")
		sb.WriteString(LabelStyle.Render("JTPCK for Claude Code, Codex, and Gemini CLI now installed!"))
		sb.WriteString("\n\n")
		if len(m.actions) > 0 {
			sb.WriteString(LabelStyle.Render("Actions:"))
			sb.WriteString("\n")
			for _, action := range m.actions {
				sb.WriteString(" - " + action + "\n")
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
		sb.WriteString(HelpStyle.Render("Restart your terminal or run: ") + CodeStyle.Render("source ~/"+m.shellConfig))
		sb.WriteString("\n")
		sb.WriteString(HelpStyle.Render("Backup: ~/" + m.shellConfig + ".jtpck-backup"))
	} else {
		// Manual installation - show copy/paste instructions
		sb.WriteString(LabelStyle.Render("To activate telemetry, add to your ~/" + m.shellConfig + ":"))
		sb.WriteString("\n\n")
		sb.WriteString(CodeStyle.Render(m.aliasCommands))
		sb.WriteString("\n\n")

		// Quick copy command
		sb.WriteString(LabelStyle.Render("Quick setup command:"))
		sb.WriteString("\n")
		copyCmd := "echo '" + strings.TrimSpace(m.aliasCommands) + "' >> ~/" + m.shellConfig
		sb.WriteString(CodeStyle.Render(copyCmd))
	}

	sb.WriteString("\n\n")
	sb.WriteString(HelpStyle.Render("Press any key to exit"))
	sb.WriteString("\n")

	return BoxStyle.Render(sb.String())
}

func (m SuccessModel) Done() bool {
	return m.done
}
