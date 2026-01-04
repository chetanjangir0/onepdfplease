package listfiles

// TODO
// fix "no items" style
// user borders
// add output file defaults and user able to edit it eg: output file: /home/user/Downloads/merge.pdf
import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/filepicker"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type file struct {
	path string
}

func (i file) FilterValue() string { return "" }

type Model struct {
	files       list.Model
	keys        keyMap
	help        help.Model
	filePicker  filepicker.Model
	pickingFile bool
	Title       string
}

func NewModel() Model {
	items := []list.Item{}

	const defaultWidth = 20
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "PDF Tool"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false) // instead using custom help menu
	l.Styles.NoItems = l.Styles.NoItems.PaddingLeft(l.Styles.TitleBar.GetPaddingLeft())

	fp := filepicker.NewModel()

	return Model{
		files:      l,
		keys:       keys,
		help:       help.New(),
		filePicker: fp,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.files.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.add):
			m.pickingFile = true
			m.filePicker = filepicker.NewModel()
			return m, m.filePicker.Init()
		case key.Matches(msg, m.keys.remove):
			m.files.RemoveItem(m.files.GlobalIndex())
			m.files.CursorUp()
			return m, nil
		case key.Matches(msg, m.keys.merge):
			log.Println("merging PDFs")
		case key.Matches(msg, m.keys.save):
			log.Println("saving PDFs")
		case key.Matches(msg, m.keys.help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.shiftDown):
			curIdx := m.files.GlobalIndex()
			m.swapItems(curIdx, curIdx+1)
			m.files.CursorDown()
			return m, nil
		case key.Matches(msg, m.keys.shiftUp):
			curIdx := m.files.GlobalIndex()
			m.swapItems(curIdx, curIdx-1)
			m.files.CursorUp()
			return m, nil
		}
	case types.QuitFilePickerMsg:
		for _, path := range msg.Paths {
			m.files.InsertItem(len(m.files.Items()), file{path: path})
		}
		m.pickingFile = false
	}

	var cmd tea.Cmd

	if m.pickingFile {
		m.filePicker, cmd = m.filePicker.Update(msg)
		return m, cmd
	}

	m.files, cmd = m.files.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	filesView := m.files.View()
	helpView := helpStyle.Render(m.help.View(m.keys))

	if m.pickingFile {
		return m.filePicker.View()
	}

	return "\n" + lipgloss.JoinVertical(
		lipgloss.Left,
		filesView,
		helpView,
	)
}

func (m *Model) swapItems(idx1, idx2 int) {
	if min(idx1, idx2) < 0 || max(idx1, idx2) >= len(m.files.Items()) {
		return
	}
	item1 := m.files.Items()[idx1]
	item2 := m.files.Items()[idx2]
	m.files.SetItem(idx1, item2)
	m.files.SetItem(idx2, item1)
}
