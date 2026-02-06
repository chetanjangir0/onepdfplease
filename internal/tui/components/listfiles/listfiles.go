package listfiles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/filepicker"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/keys"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

type file struct {
	path string
}

func (i file) FilterValue() string { return "" }

type Model struct {
	files       list.Model
	keys        keys.FileListKeymap
	help        help.Model
	filePicker  filepicker.Model
	PickingFile bool
	Title       string
	ctx         *context.ProgramContext
}

func (m *Model) SetTitle(title string) {
	m.files.Title = title
}

func (m *Model) GetFilePaths() []string {
	items := m.files.Items()
	paths := make([]string, 0, len(items))
	for _, item := range items {
		if f, ok := item.(file); ok {
			paths = append(paths, f.path)
		}
	}
	return paths
}

func (m *Model) SetAllowedFileTypes(types []string) {
	m.filePicker.SetAllowedTypes(types)
}

func NewModel(ctx *context.ProgramContext) Model {
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
	l.KeyMap.Quit.SetEnabled(false)

	fp := filepicker.NewModel(ctx)

	return Model{
		files:      l,
		keys:       keys.FileListKeys,
		help:       help.New(),
		filePicker: fp,
		ctx:        ctx,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.files.SetWidth(msg.Width)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Add):
			m.PickingFile = true
			m.filePicker.ClearSelected()
			return m, m.filePicker.Init()
		case key.Matches(msg, m.keys.Remove):
			m.files.RemoveItem(m.files.GlobalIndex())
			if m.files.Cursor() >= len(m.files.Items()) {
				m.files.CursorUp()
			}
			return m, nil
		case key.Matches(msg, m.keys.ShiftDown):
			curIdx := m.files.GlobalIndex()
			m.swapItems(curIdx, curIdx+1)
			m.files.CursorDown()
			return m, nil
		case key.Matches(msg, m.keys.ShiftUp):
			curIdx := m.files.GlobalIndex()
			m.swapItems(curIdx, curIdx-1)
			m.files.CursorUp()
			return m, nil
		case key.Matches(msg, m.keys.QuitFilepicker):
			if !m.PickingFile {
				return m, nil
			}
			for _, path := range m.filePicker.SelectedFiles {
				m.files.InsertItem(len(m.files.Items()), file{path: path})
			}
			m.PickingFile = false
			return m, nil
		}
	}

	var cmd tea.Cmd

	if m.PickingFile {
		m.filePicker, cmd = m.filePicker.Update(msg)
		return m, cmd
	}

	m.files, cmd = m.files.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	filesView := m.files.View()

	if m.PickingFile {
		return m.filePicker.View()
	}

	return "\n" + filesView
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
