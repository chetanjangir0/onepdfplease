package userinputs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
)

type Model struct {
	FocusIndex int
	Inputs     []textinput.Model
	CursorMode cursor.Mode
	ButtonText string
	Disabled   map[int]bool
	BoolInput  map[int]bool
}

type Field struct {
	Placeholder string
	Prompt      string
	IsBoolType  bool
}

func (m *Model) EnableInput(indxes []int) {
	for _, inputIdx := range indxes {
		if m.Disabled[inputIdx] {
			m.Disabled[inputIdx] = false
		}
	}
}

func (m *Model) DisableInput(indxes []int) {
	for _, inputIdx := range indxes {
		if !m.Disabled[inputIdx] {
			m.Disabled[inputIdx] = true
		}
	}
}

func (m Model) GetInputValues() []string {
	values := make([]string, len(m.Inputs))

	for i, Inp := range m.Inputs {
		values[i] = Inp.Value()
	}
	return values
}

func NewModel(fields []Field) Model {
	m := Model{
		Inputs:     make([]textinput.Model, len(fields)),
		CursorMode: cursor.CursorStatic,
		ButtonText: "Submit",
		Disabled:   make(map[int]bool),
		BoolInput:  make(map[int]bool),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64
		t.Width = 20
		if !fields[i].IsBoolType {
			t.Placeholder = fields[i].Placeholder
		} else {
			t.SetValue("no")
			m.BoolInput[i] = true
		}
		t.Prompt = fields[i].Prompt
		t.Cursor.SetMode(m.CursorMode)
		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		m.Inputs[i] = t
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.FocusIndex == len(m.Inputs) {
				return m, func() tea.Msg { return messages.OutputButtonClicked{} }
			}
		case "left", "right", "h", "l":
			if m.BoolInput[m.FocusIndex] {
				newValue := "no"
				if m.Inputs[m.FocusIndex].Value() == "no" {
					newValue = "yes"
				}
				m.Inputs[m.FocusIndex].SetValue(newValue)

				return m, func() tea.Msg {
					return messages.BoolInputToggled{
						InputIndex: m.FocusIndex,
						Value:      newValue == "yes",
					}
				}
			}
		case "up", "ctrl+p":
			return m.moveFocus(-1)
		case "down", "ctrl+n":
			return m.moveFocus(1)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m Model) moveFocus(direction int) (Model, tea.Cmd) {
	m.FocusIndex += direction

	// skip disabled fields
	for m.FocusIndex >= 0 && m.FocusIndex < len(m.Inputs) && m.Disabled[m.FocusIndex] {
		m.FocusIndex += direction
	}
	// Wraparound
	if m.FocusIndex > len(m.Inputs) {
		m.FocusIndex = 0
	} else if m.FocusIndex < 0 {
		m.FocusIndex = len(m.Inputs)
	}

	// Update focus styling
	cmds := make([]tea.Cmd, 0, len(m.Inputs))
	for i := range m.Inputs {
		if i == m.FocusIndex {
			cmds = append(cmds, m.Inputs[i].Focus())
			m.Inputs[i].PromptStyle = focusedStyle
			m.Inputs[i].TextStyle = focusedStyle
		} else {
			m.Inputs[i].Blur()
			m.Inputs[i].PromptStyle = noStyle
			m.Inputs[i].TextStyle = noStyle
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text Inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		if !m.BoolInput[i] {
			m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
		}
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteRune('\n')

	for i := range m.Inputs {
		if !m.Disabled[i] {
			b.WriteString(m.Inputs[i].View())
		}
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := fmt.Sprintf("[ %s ]", blurredStyle.Render(m.ButtonText))
	if m.FocusIndex == len(m.Inputs) {
		button = fmt.Sprintf("[ %s ]", focusedStyle.Render(m.ButtonText))
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)
	return style.DefaultStyle.MarginLeftStyle.Render(b.String())
}
