package filepicker

// TODO
// add keymaps and show keymaps in help menu
// space to toggle file selection and enter to end filepicker
// add file deletions maybe give the user an undo button
// add keys component from bubbles
// add pagination in selected items too
// don't use truncation in error view

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
)

type Model struct {
	filepicker    filepicker.Model
	SelectedFiles []string
	err           error
	ctx           *context.ProgramContext
	height        int
}

func NewModel(ctx *context.ProgramContext) Model {
	height := 20

	fp := filepicker.New()
	fp.AllowedTypes = []string{".pdf"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.SetHeight(height)
	fp.ShowPermissions = false
	// fp.KeyMap.Select = key.NewBinding(
	// 	key.WithKeys(" "),
	// 	key.WithHelp("space", "select"),
	// )
	return Model{
		filepicker: fp,
		ctx:        ctx,
		height:     height,
	}
}

func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+y":
			return m, func() tea.Msg {
				return messages.QuitFilePicker{Paths: m.SelectedFiles} // TODO: use reference here
			}
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.SelectedFiles = append(m.SelectedFiles, path)
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		// m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
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
	view.WriteString("\n  ")
	view.WriteString("Pick files:")
	view.WriteString("\n\n" + m.filepicker.View() + "\n")

	if m.err != nil {
		view.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	}
	return view.String()
}

func (m Model) selectedView() string {
	var view strings.Builder
	view.WriteString("\n  ")
	view.WriteString("Selected files: \n")
	view.WriteString("\n")
	for i, f := range m.SelectedFiles {
		// only show the last m.height files when files are too many
		if m.height <= len(m.SelectedFiles) && i < len(m.SelectedFiles)-m.height {
			continue
		}
		view.WriteString(m.filepicker.Styles.Selected.PaddingLeft(2).Render(filepath.Base(f)) + "\n")
	}
	view.WriteString("\n")
	return view.String()
}
