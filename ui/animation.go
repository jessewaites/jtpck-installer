package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AnimationModel struct {
	frame  int
	done   bool
	height int
	width  int
	frames []string
	delays []time.Duration
	useGif bool
}

type AnimationDoneMsg struct{}

func NewAnimationModel() AnimationModel {
	// Try to load the GIF at startup; fall back to baked ASCII frames.
	width := chooseGifWidth(0)
	gifFrames, gifDelays, ok := loadGifFrames(width)

	frameSet := frames
	delaySet := make([]time.Duration, len(frameSet))
	for i := range delaySet {
		delaySet[i] = 100 * time.Millisecond
	}

	useGif := false
	if ok && len(gifFrames) > 0 {
		frameSet = gifFrames
		delaySet = gifDelays
		useGif = true
	}

	return AnimationModel{
		frame:  0,
		frames: frameSet,
		delays: delaySet,
		useGif: useGif,
	}
}

func (m AnimationModel) Init() tea.Cmd {
	return tea.Tick(m.currentDelay(), func(t time.Time) tea.Msg {
		return AnimationDoneMsg{}
	})
}

func (m AnimationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	totalFrames := len(m.frames)
	if totalFrames == 0 {
		m.done = true
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		if gifFrames, gifDelays, ok := loadGifFrames(chooseGifWidth(msg.Width)); ok {
			m.frames = gifFrames
			m.delays = gifDelays
			m.useGif = true
			if m.frame >= len(m.frames) {
				m.frame = 0
			}
		}
		return m, nil

	case tea.KeyMsg:
		// Allow skipping animation with any key
		m.done = true
		return m, tea.Quit

	case AnimationDoneMsg:
		m.frame++
		if m.frame >= totalFrames {
			m.done = true
			return m, tea.Quit
		}
		return m, tea.Tick(m.currentDelay(), func(t time.Time) tea.Msg {
			return AnimationDoneMsg{}
		})
	}

	return m, nil
}

func (m AnimationModel) View() string {
	if m.done || len(m.frames) == 0 {
		return ""
	}

	// Use imported animation frames
	totalFrames := len(m.frames)
	currentFrame := m.frames[m.frame%totalFrames]
	if !m.useGif {
		currentFrame = strings.ReplaceAll(currentFrame, "@", " ")
	}

	frameLines := strings.Split(currentFrame, "\n")

	const extraLines = 6 // title + progress + help + spacing
	targetHeight := m.height
	if targetHeight <= 0 {
		targetHeight = len(frameLines) + extraLines
	}
	if !m.useGif {
		maxFrameLines := targetHeight - extraLines
		if maxFrameLines < 0 {
			maxFrameLines = 0
		}
		if len(frameLines) > maxFrameLines && maxFrameLines > 0 {
			step := (len(frameLines) + maxFrameLines - 1) / maxFrameLines
			if step < 1 {
				step = 1
			}
			compressed := make([]string, 0, maxFrameLines)
			for i := 0; i < len(frameLines) && len(compressed) < maxFrameLines; i += step {
				compressed = append(compressed, frameLines[i])
			}
			frameLines = compressed
		}
	}

	var renderedFrame string
	containerWidth := 90
	if m.width > 0 {
		containerWidth = m.width
	}

	frameBase := lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Width(containerWidth).
		Align(lipgloss.Center)

	if m.useGif {
		renderedFrame = frameBase.Render(strings.Join(frameLines, "\n"))
	} else {
		renderedFrame = frameBase.Foreground(BrandRed).Render(strings.Join(frameLines, "\n"))
	}

	// Title styling
	title := lipgloss.NewStyle().
		Foreground(BrandRed).
		Bold(true).
		Render("═══ Installing JTPCK for AI Telemetry  ═══")
	title = frameBase.Render(title)

	progressBar := generateProgressBar(m.frame+1, totalFrames)
	progressBar = frameBase.Render(progressBar)

	// Combine all elements
	var sb strings.Builder

	// Calculate vertical padding to center content
	contentLines := len(frameLines) + 4 // frame + title + progress + help
	if m.height > contentLines {
		topPadding := (m.height - contentLines) / 3 // More space at bottom
		bottomPadding := (m.height - contentLines) - topPadding

		// Add top padding
		for i := 0; i < topPadding; i++ {
			sb.WriteString("\n")
		}

		sb.WriteString(renderedFrame)
		sb.WriteString("\n")
		// Add blank lines above title
		sb.WriteString(frameBase.Render(""))
		sb.WriteString("\n")
		sb.WriteString(frameBase.Render(""))
		sb.WriteString("\n")
		sb.WriteString(title)
		sb.WriteString("\n")
		sb.WriteString(progressBar)
		sb.WriteString("\n")
		sb.WriteString(frameBase.Render(HelpStyle.Render("Press any key to skip")))

		// Add bottom padding
		for i := 0; i < bottomPadding; i++ {
			sb.WriteString("\n")
		}
	} else {
		// No padding if terminal too small
		sb.WriteString(renderedFrame)
		sb.WriteString("\n")
		// Add blank lines above title
		sb.WriteString(frameBase.Render(""))
		sb.WriteString("\n")
		sb.WriteString(frameBase.Render(""))
		sb.WriteString("\n")
		sb.WriteString(title)
		sb.WriteString("\n")
		sb.WriteString(progressBar)
		sb.WriteString("\n")
		sb.WriteString(frameBase.Render(HelpStyle.Render("Press any key to skip")))
	}

	return sb.String()
}

func generateProgressBar(current, total int) string {
	if total <= 0 {
		return ""
	}
	if current > total {
		current = total
	}

	percentage := (current * 100) / total
	filled := (current * 40) / total

	bar := strings.Repeat("█", filled) + strings.Repeat("░", 40-filled)
	return lipgloss.NewStyle().
		Foreground(BrandRed).
		Render(fmt.Sprintf("%s %d%%", bar, percentage))
}

func (m AnimationModel) Done() bool {
	return m.done
}

func (m AnimationModel) currentDelay() time.Duration {
	if m.useGif && len(m.delays) > 0 {
		// GIF delays are in 10ms units; convert to duration and slow down a bit for readability.
		d := m.delays[m.frame%len(m.delays)] * 20
		if d < 33*time.Millisecond {
			d = 33 * time.Millisecond
		}
		return d
	}
	return 100 * time.Millisecond
}
