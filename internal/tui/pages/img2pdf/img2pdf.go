package img2pdf

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/listfiles"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/userinputs"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
	"github.com/chetanjangir0/onepdfplease/internal/tui/utils"
)

const (
	outputFileIdx = iota
	mergeIntoOneIdx
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
		outputPlaceholder: "./converted.pdf",
	}
	lf := listfiles.NewModel(ctx)
	lf.SetTitle("Choose Order")
	lf.SetAllowedFileTypes([]string{".png", ".jpg", ".jpeg", ".tif", ".webp"})

	outputFields := make([]userinputs.Field, 2)
	outputFields[outputFileIdx] = userinputs.Field{
		Placeholder: m.outputPlaceholder,
		Prompt:      "Output File: ",
	}
	outputFields[mergeIntoOneIdx] = userinputs.Field{
		Placeholder: m.outputPlaceholder,
		Prompt:      "convert and Merge into one file?: ",
		IsBoolType:  true,
	}

	op := userinputs.NewModel(outputFields)
	op.ButtonText = "Convert and Save"

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
		switch msg.String() {
		case "tab": // switch focus
			m.focusIndex = (m.focusIndex + 1) % 2
			return m, nil
		case "shift+tab":
			m.focusIndex = (m.focusIndex - 1 + 2) % 2
			return m, nil
		case "esc":
			return m, func() tea.Msg {
				return messages.Navigate{Page: types.MenuPage}
			}
		}
	case messages.OutputButtonClicked:
		outFile := m.outputPlaceholder
		mergeIntoOne := false

		userValues := m.outputPicker.GetInputValues()
		if len(userValues) > outputFileIdx && len(userValues[outputFileIdx]) != 0 {
			outFile = userValues[outputFileIdx]
		}
		if len(userValues) > mergeIntoOneIdx && userValues[mergeIntoOneIdx] == "yes" {
			mergeIntoOne = true
		}
		return m, utils.Img2Pdf(m.fileList.GetFilePaths(), outFile, mergeIntoOne)
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
