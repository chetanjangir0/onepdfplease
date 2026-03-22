package merge

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/pdf"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/listfiles"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/userinputs"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/keys"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

const (
	outputFileIdx = iota
)

type Model struct {
	focusIndex        int // 0 for fileList 1 for outputPicker
	fileList          listfiles.Model
	outputPicker      userinputs.Model
	ctx               *context.ProgramContext
	outputPlaceholder string
}

func NewModel(ctx *context.ProgramContext) Model {
	m := Model{
		outputPlaceholder: "./merged.pdf",
	}
	lf := listfiles.NewModel(ctx)
	lf.SetTitle("Choose Order")

	outputFields := make([]userinputs.Field, 1)
	outputFields[outputFileIdx] = userinputs.Field{
		Placeholder: m.outputPlaceholder,
		Prompt:      "Output File: ",
	}

	op := userinputs.NewModel(outputFields)
	op.ButtonText = "Merge and Save"

	m.fileList = lf
	m.outputPicker = op
	m.ctx = ctx
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

    // If file list is picking, give it complete control
    if m.fileList.PickingFile {
        m.fileList, cmd = m.fileList.Update(msg)
        return m, cmd
    }

	switch m.focusIndex {
	case 0:
		m.fileList, cmd = m.fileList.Update(msg)
	case 1:
		m.outputPicker, cmd = m.outputPicker.Update(msg)
	}

	if cmd != nil {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.NxtTab):
			m.focusIndex = (m.focusIndex + 1) % 2
			return m, nil
		case key.Matches(msg, keys.Keys.PrevTab):
			m.focusIndex = (m.focusIndex - 1 + 2) % 2
			return m, nil
		case key.Matches(msg, keys.Keys.Back):
			return m, func() tea.Msg {
				return messages.Navigate{Page: types.MenuPage}
			}
		}
	case messages.OutputButtonClicked:
		outFile := m.outputPlaceholder
		userValues := m.outputPicker.GetInputValues()
		if len(userValues) > outputFileIdx && len(userValues[outputFileIdx]) != 0 {
			outFile = userValues[outputFileIdx]
		}
		return m, pdf.Merge(m.fileList.GetFilePaths(), outFile, m.ctx)
	}


	return m, cmd
}

func (m Model) View() string {
	if m.fileList.PickingFile {
		return m.fileList.View()
	}
	var fileListStyle, outputPickerStyle lipgloss.Style
	if m.focusIndex == 0 {
		fileListStyle = style.DefaultStyle.FocusedBorder
		outputPickerStyle = style.DefaultStyle.BlurredBorder
	} else {
		fileListStyle = style.DefaultStyle.BlurredBorder
		outputPickerStyle = style.DefaultStyle.FocusedBorder
	}

	return style.RenderTwoFullRows(
		m.ctx.TermWidth,
		m.ctx.MainContentHeight,
		fileListStyle,
		outputPickerStyle,
		m.fileList.View(),
		m.outputPicker.View(),
	)
}
