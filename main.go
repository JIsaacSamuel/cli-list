package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error

type model struct {
	textarea textarea.Model
	fstyle   lipgloss.Style
	err      error
	width    int
	height   int
}

func initialModel() model {
	ti := textarea.New()
	ti.Placeholder = "Your notes here..."
	ti.SetWidth(80)
	ti.Focus()

	appstyle := lipgloss.NewStyle().
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	return model{
		textarea: ti,
		fstyle:   appstyle,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return fmt.Sprintf("Loading program...")
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		fmt.Sprintf(
			"CLI notepad.\n\n%s\n\n%s",
			m.fstyle.Render(m.textarea.View()),
			"(ctrl+c to quit)",
		)+"\n\n",
	)
}
