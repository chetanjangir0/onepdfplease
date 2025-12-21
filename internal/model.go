package internal

import tea "github.com/charmbracelet/bubbletea"

type MenuState int

const (
	Tools MenuState = iota
	Picker
	Merge
	Split
)

type model struct {
	cursor        int
	CurrentMenu   MenuState
	ToolsMenu     []string
	SelectedTool  string
	SelectedFiles []string
	status        string
	width         int
	height        int
}

func InitialModel() model {

	return model{
		cursor:      0,
		CurrentMenu: Tools,
		ToolsMenu:   []string{"Merge PDFs", "Split PDF"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "j", "down":
			if m.cursor < m.itemCount()-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter":
			switch m.CurrentMenu {
			case Tools:

				switch m.ToolsMenu[m.cursor] {
				case "Merge PDFs":
					m.CurrentMenu = Merge
					m.cursor = 0
					return m, nil
				case "Split PDF":
					m.CurrentMenu = Split
					m.cursor = 0
					return m, nil
				}
				m.cursor = 0 // reset cursor pos

			case Merge:
				m.status = "Merge pdfs"
				return m, nil
			case Split:
				m.status = "Split pdf"
				return m, nil

			}
		case "b", "esc":
			m.CurrentMenu = Tools
			m.cursor = 0
			m.status = ""
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.status
}

func (m model) itemCount() int {
	switch m.CurrentMenu {
	case Tools:
		return len(m.ToolsMenu)
	default:
		return 0
	}
}
