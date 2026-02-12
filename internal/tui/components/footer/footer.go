package footer

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/keys"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
)

type Model struct {
	help    help.Model
	ctx     *context.ProgramContext
	spinner spinner.Model
	ShowAll bool
}

func NewModel(ctx *context.ProgramContext) Model {
	help := help.New()
	help.ShowAll = true

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m := Model{
		help:    help,
		ctx:     ctx,
		spinner: s,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	var statusStyle lipgloss.Style
	var icon string

	switch m.ctx.StatusType {
	case context.Error:
		statusStyle = style.DefaultStyle.ErrorStyle
		icon = "✗"
	case context.Success:
		statusStyle = style.DefaultStyle.SuccessStyle
		icon = "✓"
	case context.Processing:
		statusStyle = style.DefaultStyle.NeutralStyle
		icon = m.spinner.View()
	default:
		statusStyle = style.DefaultStyle.NeutralStyle
		icon = "Press ? for help"
	}

	footer := statusStyle.Render(fmt.Sprintf("%s %s", icon, m.ctx.Status))

	if m.ShowAll {
		keyMap := keys.CreateKeyMapForPage(m.ctx.CurrentPage)
		fullHelp := m.help.View(keyMap)
		footer = lipgloss.JoinVertical(lipgloss.Left, footer, fullHelp)
	}
	return lipgloss.NewStyle().MarginLeft(2).Render(footer)
}
