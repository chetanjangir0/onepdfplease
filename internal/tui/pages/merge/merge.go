package merge

// TODO:
// fix "no items" style
// user borders
// add output file defaults and user able to edit it eg: output file: /home/user/Downloads/merge.pdf
// change focus using tab
import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/filepicker"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/outputpicker"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type file struct {
	path string
}

func (i file) FilterValue() string { return "" }

type Model struct {
	files        list.Model
	choice       string
	keys         keyMap
	help         help.Model
	filePicker   filepicker.Model
	pickingFile  bool
	outputPicker outputpicker.Model
}

func NewModel() Model {
	items := []list.Item{}

	const defaultWidth = 20
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "PDF merger"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false) // instead using custom help menu
	l.Styles.NoItems = l.Styles.NoItems.PaddingLeft(l.Styles.TitleBar.GetPaddingLeft())

	fp := filepicker.NewModel()

	outputFields := []outputpicker.Field{
		{
			Placeholder: "./merged.pdf",
			Prefix:      "Output File",
		},
	}
	op := outputpicker.NewModel(outputFields)
	return Model{
		files:        l,
		keys:         keys,
		help:         help.New(),
		filePicker:   fp,
		outputPicker: op,
	}
}

func (m Model) Init() tea.Cmd {
	m.outputPicker.Init()
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
		case key.Matches(msg, m.keys.tab):
			m.outputPicker.Focused = !m.outputPicker.Focused
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

	if m.outputPicker.Focused {
		m.outputPicker, cmd = m.outputPicker.Update(msg)
		return m, cmd
	}

	m.files, cmd = m.files.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Initiating", m.choice))
	}

	filesView := m.files.View()
	helpView := helpStyle.Render(m.help.View(m.keys))
	outputPickerView := m.outputPicker.View()

	if m.pickingFile {
		return m.filePicker.View()
	}

	return "\n" + lipgloss.JoinVertical(
		lipgloss.Left,
		filesView,
		outputPickerView,
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
