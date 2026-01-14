package encrypt

// TODO
// add option to do inplace encryption

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
	passwdIdx = iota
	pathIdx
	prefixIdx
)

type Model struct {
	focusIndex        int // 0 for fileList 1 for outputPicker
	fileList          listfiles.Model
	userInputs        userinputs.Model
	ctx               *context.ProgramContext
	pathPlaceholder   string
	prefixPlaceholder string
}

func NewModel(ctx *context.ProgramContext) Model {
	m := Model{
		pathPlaceholder:   "./",
		prefixPlaceholder: "encrypted_",
	}
	lf := listfiles.NewModel(ctx)
	lf.SetTitle("Choose Files")

	inputFields := make([]userinputs.Field, 3)
	inputFields[passwdIdx] = userinputs.Field{
		Placeholder: "",
		Prompt:      "Password: ",
	}
	inputFields[pathIdx] = userinputs.Field{
		Placeholder: m.pathPlaceholder,
		Prompt:      "Output Path: ",
	}
	inputFields[prefixIdx] = userinputs.Field{
		Placeholder: m.prefixPlaceholder,
		Prompt:      "Output File Prefix: ",
	}

	m.userInputs = userinputs.NewModel(inputFields)
	m.userInputs.ButtonText = "Encrypt and Save"
	m.fileList = lf
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
		// TODO: use keymaps
		case "tab": // switch focus
			m.focusIndex = (m.focusIndex + 1) % 2
			return m, nil
		case "shift+tab":
			m.focusIndex = (m.focusIndex - 1 + 2) % 2
			return m, nil
		}
	case messages.OutputButtonClicked:
		outPath := m.pathPlaceholder
		outPrefix := m.prefixPlaceholder
		password := ""

		userValues := m.userInputs.GetInputValues()
		if len(userValues) > passwdIdx { // the encrypt function will handle 0 length password
			password = userValues[passwdIdx]
		}
		if len(userValues) > pathIdx && len(userValues[pathIdx]) != 0 {
			outPath = userValues[pathIdx]
		}
		if len(userValues) > prefixIdx && len(userValues[prefixIdx]) != 0 {
			outPrefix = userValues[prefixIdx]
		}
		return m, utils.Encrypt(m.fileList.GetFilePaths(), password, outPath, outPrefix)
	}

	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		m.fileList, cmd = m.fileList.Update(msg)
	case 1:
		m.userInputs, cmd = m.userInputs.Update(msg)
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
