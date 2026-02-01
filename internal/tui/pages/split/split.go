package split

// TODO:
// add verification to ranges of pages
// set option to extract all pages to yes by default

import (
	"strings"

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
	exteractAllIdx = iota
	mergeIntoOneIdx
	pathIdx
	prefixIdx
	selectedPagesIdx
)

type Model struct {
	focusIndex               int // 0 for fileList 1 for outputPicker
	fileList                 listfiles.Model
	userInputs               userinputs.Model
	ctx                      *context.ProgramContext
	pathPlaceholder          string
	prefixPlaceholder        string
	selectedPagesPlaceholder string
}

func NewModel(ctx *context.ProgramContext) Model {
	m := Model{
		pathPlaceholder:          "./",
		prefixPlaceholder:        "split_",
		selectedPagesPlaceholder: `1,3,even,4-10,2-`,
	}
	lf := listfiles.NewModel(ctx)
	lf.SetTitle("Choose Files")

	inputFields := make([]userinputs.Field, 5)
	inputFields[exteractAllIdx] = userinputs.Field{
		Placeholder: "",
		Prompt:      "Extract all pages into seperate file?: ",
		IsBoolType:  true,
	}
	inputFields[mergeIntoOneIdx] = userinputs.Field{
		Placeholder: "",
		Prompt:      "Split and Merge into one file?: ",
		IsBoolType:  true,
	}
	inputFields[pathIdx] = userinputs.Field{
		Placeholder: m.pathPlaceholder,
		Prompt:      "Output Path: ",
	}
	inputFields[prefixIdx] = userinputs.Field{
		Placeholder: m.prefixPlaceholder,
		Prompt:      "Output File Prefix: ",
	}
	inputFields[selectedPagesIdx] = userinputs.Field{
		Placeholder: m.selectedPagesPlaceholder,
		Prompt:      "Select Pages: ",
	}

	m.userInputs = userinputs.NewModel(inputFields)
	m.userInputs.ButtonText = "Split and Save"
	m.fileList = lf
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
		m.userInputs, cmd = m.userInputs.Update(msg)
	}

	if cmd != nil {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// TODO: use keymaps
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
	case messages.BoolInputToggled:
		if msg.InputIndex == exteractAllIdx {
			if msg.Value {
				m.userInputs.DisableInput([]int{mergeIntoOneIdx, pathIdx, prefixIdx, selectedPagesIdx})
			} else {
				m.userInputs.EnableInput([]int{mergeIntoOneIdx, pathIdx, prefixIdx, selectedPagesIdx})
			}
		}
		return m, nil
	case messages.OutputButtonClicked:
		outPath := m.pathPlaceholder
		outPrefix := m.prefixPlaceholder
		selectedPages := ""
		mergeIntoOne := false
		extractAllPages := false

		userValues := m.userInputs.GetInputValues()
		if len(userValues) > exteractAllIdx && userValues[exteractAllIdx] == "yes" {
			extractAllPages = true
		}
		if len(userValues) > mergeIntoOneIdx && userValues[mergeIntoOneIdx] == "yes" {
			mergeIntoOne = true
		}
		if len(userValues) > pathIdx && len(userValues[pathIdx]) != 0 {
			outPath = userValues[pathIdx]
		}
		if len(userValues) > prefixIdx && len(userValues[prefixIdx]) != 0 {
			outPrefix = userValues[prefixIdx]
		}
		if len(userValues) > selectedPagesIdx && len(userValues[selectedPagesIdx]) != 0 {
			selectedPages = userValues[selectedPagesIdx]
		}

		return m, utils.Split(
			m.fileList.GetFilePaths(),
			outPath,
			outPrefix,
			strings.Split(selectedPages, ","),
			mergeIntoOne,
			extractAllPages,
			m.ctx,
		)
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
