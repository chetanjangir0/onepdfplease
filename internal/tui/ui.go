package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/menu"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/merge"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

type model struct {
	quitting    bool
	currentPage types.Page
	ctx         *context.ProgramContext

	// each page has its own model
	menuModel  menu.Model
	mergeModel merge.Model
}

func InitialModel() model {
	m := model{
		currentPage: types.MenuPage,
	}
	m.ctx = &context.ProgramContext{}
	m.menuModel = menu.NewModel(m.ctx)
	m.mergeModel = merge.NewModel(m.ctx)
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case messages.Navigate:
		m.currentPage = msg.Page
		return m, nil
	}

	var cmd tea.Cmd
	switch m.currentPage {
	case types.MenuPage:
		m.menuModel, cmd = m.menuModel.Update(msg)
	case types.MergePage:
		m.mergeModel, cmd = m.mergeModel.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	switch m.currentPage {
	case types.MenuPage:
		return m.menuModel.View()
	case types.MergePage:
		return m.mergeModel.View()
	}
	return ""

}

func (m *model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	log.Println("window size changed", "width", msg.Width, "height", msg.Height)
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
}
