package ui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var uuidRegex = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

type InputModel struct {
	textInput    textinput.Model
	err          string
	submitted    bool
	currentValue string
	userID       string
}

func NewInputModel(currentValue string) InputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter your user ID (e.g., UUID)"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	if currentValue != "" {
		ti.SetValue(currentValue)
	}

	return InputModel{
		textInput:    ti,
		currentValue: currentValue,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Validate input
			value := strings.TrimSpace(strings.ToLower(m.textInput.Value()))
			if value == "" {
				m.err = "User ID cannot be empty"
				return m, nil
			}
			if !uuidRegex.MatchString(value) {
				m.err = "User ID must be a valid UUID (e.g., 12345678-1234-1234-1234-123456789abc)"
				return m, nil
			}
			m.userID = value
			m.submitted = true
			m.err = ""
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	if m.submitted {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(TitleStyle.Render("JTPCK Telemetry Setup"))
	sb.WriteString("\n\n")

	// Show current value if reconfiguring
	if m.currentValue != "" {
		sb.WriteString(LabelStyle.Render("Current User ID: "))
		sb.WriteString(CodeStyle.Render(m.currentValue))
		sb.WriteString("\n\n")
	}

	sb.WriteString(LabelStyle.Render("Enter your user ID:"))
	sb.WriteString("\n\n")
	sb.WriteString(m.textInput.View())
	sb.WriteString("\n\n")

	if m.err != "" {
		sb.WriteString(WarningStyle.Render("⚠ " + m.err))
		sb.WriteString("\n\n")
	}

	sb.WriteString(HelpStyle.Render("Press Enter to continue • Ctrl+C to cancel"))
	sb.WriteString("\n")

	return BoxStyle.Render(sb.String())
}

func (m InputModel) Done() bool {
	return m.submitted
}

func (m InputModel) GetUserID() string {
	return m.userID
}
