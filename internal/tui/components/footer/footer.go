package footer

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/keys"
)

type Model struct {
	help    help.Model
	ctx     *context.ProgramContext
	ShowAll bool
}

func NewModel(ctx *context.ProgramContext) Model {
	help := help.New()
	help.ShowAll = true
	m := Model{
		help: help,
	}
	return m
}

func (m Model) View() string {
	footer := "This is the footer"

	if m.ShowAll {
		fullHelp := m.help.View(keys.Keys)
		return lipgloss.JoinVertical(lipgloss.Top, footer, fullHelp)
	}
	return footer
}
