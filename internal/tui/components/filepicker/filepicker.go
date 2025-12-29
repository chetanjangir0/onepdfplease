package filepicker

// TODO
// make two columns for selected files and file View
// add keymaps and show keymaps in help menu
// space to toggle file selection and enter to end filepicker
// add file deletions maybe give the user an undo button
// add keys component from bubbles
// add swap mechanism

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

type Model struct {
	filepicker    filepicker.Model
	SelectedFiles []string
	err           error
}

func NewModel() Model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".pdf"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.SetHeight(20)
	// fp.KeyMap.Select = key.NewBinding(
	// 	key.WithKeys(" "),
	// 	key.WithHelp("space", "select"),
	// )
	return Model{
		filepicker: fp,
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
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+y":
			return m, func() tea.Msg {
				return types.QuitFilePickerMsg{Paths: m.SelectedFiles} // TODO: use reference here
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
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if len(m.SelectedFiles) == 0 {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected files: \n")
		for _, f := range m.SelectedFiles {
			s.WriteString(m.filepicker.Styles.Selected.Render(f) + "\n")
		}
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}
