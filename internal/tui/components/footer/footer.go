package footer

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/keys"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
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
		ctx:  ctx,
	}
	return m
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
	default:
		statusStyle = style.DefaultStyle.NeutralStyle
		icon = "i"
	}

	footer := statusStyle.Render(fmt.Sprintf("%s %s", icon, m.ctx.Status))

	if m.ShowAll {
		keyMap := keys.CreateKeyMapForPage(m.ctx.CurrentPage)
		fullHelp := m.help.View(keyMap)
		return lipgloss.JoinVertical(lipgloss.Top, footer, fullHelp)
	}
	return footer
}

func (m *Model) ShowError(msg string) {
	m.ctx.Status = msg
	m.ctx.StatusType = context.Error
}

func (m *Model) ShowSuccess(msg string) {
	m.ctx.Status = msg
	m.ctx.StatusType = context.Success
}

func (m *Model) ClearStatus() {
	m.ctx.Status = ""
	m.ctx.StatusType = context.None
}
