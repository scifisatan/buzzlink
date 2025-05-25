package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type uploadModel struct {
	spinner spinner.Model
	done    bool
	err     error
	link    string
}

func (m uploadModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m uploadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		if m.done {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case uploadDoneMsg:
		m.done = true
		m.link = msg.link
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

type uploadDoneMsg struct {
	link string
	err  error
}

func (m uploadModel) View() string {
	if m.done {
		if m.err != nil {
			return fmt.Sprintf("%s Upload failed: %v\n", IconError, m.err)
		}
		return fmt.Sprintf("%s %sUploaded successfully!%s\n", IconSuccess, ColorGreen, ColorReset)
	}
	return fmt.Sprintf("%s %s Uploading...", IconUpload, m.spinner.View())
}

// runUploadWithSpinner runs the upload with a spinner animation
func runUploadWithSpinner(zippedPath, note string) (string, error) {
	spin := spinner.New()
	spin.Spinner = spinner.Dot                               // Use the basic dot spinner
	spin.Style = spin.Style.Foreground(lipgloss.Color("36")) // Cyan/blue
	m := uploadModel{spinner: spin}

	done := make(chan uploadDoneMsg, 1)
	// Start upload in goroutine
	go func() {
		l, e := UploadFile(zippedPath, note)
		done <- uploadDoneMsg{link: l, err: e}
	}()

	// Bubble Tea program with custom update to receive uploadDoneMsg
	p := tea.NewProgram(m)
	var result uploadDoneMsg
	go func() {
		result = <-done
		p.Send(result)
	}()
	if _, err := p.Run(); err != nil {
		return "", fmt.Errorf("spinner error: %w", err)
	}
	return result.link, result.err
}
