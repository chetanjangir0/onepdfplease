package doc2pdf

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/listfiles"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/userinputs"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/keys"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
	"github.com/chetanjangir0/onepdfplease/internal/tui/utils"
)

const (
	pathIdx = iota
	prefixIdx
)

type Model struct {
	focusIndex        int
	fileList          listfiles.Model
	userInputs        userinputs.Model
	ctx               *context.ProgramContext
	pathPlaceholder   string
	prefixPlaceholder string
}

func NewModel(ctx *context.ProgramContext) Model {
	m := Model{
		pathPlaceholder:   "./",
		prefixPlaceholder: "converted_",
	}

	lf := listfiles.NewModel(ctx)
	lf.SetTitle("Choose Files")
	lf.SetAllowedFileTypes([]string{".doc", ".docx"})

	inputFields := make([]userinputs.Field, 2)
	inputFields[pathIdx] = userinputs.Field{
		Placeholder: m.pathPlaceholder,
		Prompt:      "Output Path: ",
	}
	inputFields[prefixIdx] = userinputs.Field{
		Placeholder: m.prefixPlaceholder,
		Prompt:      "Output File Prefix: ",
	}

	m.userInputs = userinputs.NewModel(inputFields)
	m.userInputs.ButtonText = "Convert and Save"
	m.fileList = lf
	m.ctx = ctx

	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.fileList.PickingFile {
		m.fileList, cmd = m.fileList.Update(msg)
		return m, cmd
	}

	switch m.focusIndex {
	case 0:
		m.fileList, cmd = m.fileList.Update(msg)
	case 1:
		m.userInputs, cmd = m.userInputs.Update(msg)
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
		outPath := m.pathPlaceholder
		outPrefix := m.prefixPlaceholder

		userValues := m.userInputs.GetInputValues()
		if len(userValues) > pathIdx && len(userValues[pathIdx]) != 0 {
			outPath = userValues[pathIdx]
		}
		if len(userValues) > prefixIdx && len(userValues[prefixIdx]) != 0 {
			outPrefix = userValues[prefixIdx]
		}

		return m, utils.Doc2Pdf(m.fileList.GetFilePaths(), outPath, outPrefix, m.ctx)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.fileList.PickingFile {
		return m.fileList.View()
	}

	var fileListStyle, userInputStyle lipgloss.Style
	if m.focusIndex == 0 {
		fileListStyle = style.DefaultStyle.FocusedBorder
		userInputStyle = style.DefaultStyle.BlurredBorder
	} else {
		fileListStyle = style.DefaultStyle.BlurredBorder
		userInputStyle = style.DefaultStyle.FocusedBorder
	}

	return style.RenderTwoFullRows(
		m.ctx.TermWidth,
		m.ctx.MainContentHeight,
		fileListStyle,
		userInputStyle,
		m.fileList.View(),
		m.userInputs.View(),
	)
}
