package merge

// TODO:
// change focus using tab
// fix border changing after selected files become active
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/listfiles"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/userinputs"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
	"github.com/chetanjangir0/onepdfplease/internal/tui/utils"
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab": // switch focus
			m.focusIndex = (m.focusIndex + 1) % 2
			return m, nil
		case "shift+tab":
			m.focusIndex = (m.focusIndex - 1 + 2) % 2
			return m, nil
		}
	case messages.OutputButtonClicked:
		outFile := m.outputPlaceholder
		userValues := m.outputPicker.GetInputValues()
		if len(userValues) > outputFileIdx && len(userValues[outputFileIdx]) != 0 {
			outFile = userValues[outputFileIdx]
		}
		return m, utils.Merge(m.fileList.GetFilePaths(), outFile)
	}

	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		m.fileList, cmd = m.fileList.Update(msg)
	case 1:
		m.outputPicker, cmd = m.outputPicker.Update(msg)
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
