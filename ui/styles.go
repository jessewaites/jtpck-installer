package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Brand colors
	BrandRed   = lipgloss.Color("#FF0000")
	BrandWhite = lipgloss.Color("#FFFFFF")
	BrandGray  = lipgloss.Color("#888888")

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Foreground(BrandRed).
			Bold(true).
			Padding(1, 0)

	// Box style
	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(BrandRed).
			Padding(1, 2).
			Width(60)

	// Input style
	InputStyle = lipgloss.NewStyle().
			Foreground(BrandWhite).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1)

	// Label style
	LabelStyle = lipgloss.NewStyle().
			Foreground(BrandWhite).
			Bold(true)

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(BrandGray).
			Italic(true)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	// Warning style
	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	// Code block style
	CodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Background(lipgloss.Color("#1A1A1A")).
			Padding(0, 1)
)
