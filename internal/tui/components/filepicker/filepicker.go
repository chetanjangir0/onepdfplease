package filepicker

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
)

var browseViewHeader = "\n  Pick files:\n\n"
var selectedViewHeader = "\n  Selected files:\n\n"

type Model struct {
	filepicker       filepicker.Model
	SelectedFiles    []string
	ctx              *context.ProgramContext
	maxContentHeight int
	contentHeight    int
}

func NewModel(ctx *context.ProgramContext) Model {
	const maxContentHeight = 20

	fp := filepicker.New()
	fp.AllowedTypes = []string{".pdf"}
	fp.CurrentDirectory, _ = os.Getwd()
	fp.SetHeight(maxContentHeight - 1)
	fp.AutoHeight = false
	fp.ShowPermissions = false
	return Model{
		filepicker:       fp,
		ctx:              ctx,
		maxContentHeight: maxContentHeight,
		contentHeight:    maxContentHeight,
	}
}

func (m *Model) onWindowSizeChanged() {
	headerHeight := lipgloss.Height(browseViewHeader) 
	borderHeight := 1
	usedHeight := headerHeight + m.contentHeight + 2*borderHeight
	if m.ctx.MainContentHeight < usedHeight {
		m.contentHeight = m.ctx.MainContentHeight - headerHeight - 2*borderHeight
	} else {
		m.contentHeight = m.maxContentHeight
	}
	m.filepicker.SetHeight(m.contentHeight - 1)// -1 accounts for the extra 1 height upstreams sets
}

func (m *Model) ClearSelected() {
	m.SelectedFiles = nil
}

func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m *Model) SetAllowedTypes(types []string) {
	m.filepicker.AllowedTypes = types
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged()
	case tea.KeyMsg:
		switch msg.String() {
		// case "ctrl+y":
		}
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.SelectedFiles = append(m.SelectedFiles, path)
	}

	// Did the user select a disabled file?
	if didSelect, _ := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		err := errors.New("File is not valid.")

		return m, func() tea.Msg {
			return messages.ShowError{
				Err: err,
			}
		}

	}

	return m, cmd
}

func (m Model) View() string {
	return style.RenderTwoFullCols(
		m.ctx.TermWidth,
		m.ctx.MainContentHeight,
		style.DefaultStyle.FocusedBorder,
		m.browseView(),
		m.selectedView(),
	)
}

func (m Model) browseView() string {
	var view strings.Builder
	view.WriteString(browseViewHeader)
	view.WriteString(m.filepicker.View())

	return view.String()
}

func (m Model) selectedView() string {
	var view strings.Builder
	view.WriteString(selectedViewHeader)
	for i, f := range m.SelectedFiles {
		// only show the last m.contentHeight files when files are too many
		if m.contentHeight-1 <= len(m.SelectedFiles) && i < len(m.SelectedFiles)-(m.contentHeight - 1){
			continue
		}
		view.WriteString(m.filepicker.Styles.Selected.PaddingLeft(2).Render(filepath.Base(f)) + "\n")
	}
	return view.String()
}
