package merge

// TODO:
// change focus using tab
// fix border changing after selected files become active
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/listfiles"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/outputpicker"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
)

type Model struct {
	focusIndex   int // 0 for fileList 1 for outputPicker
	fileList     listfiles.Model
	outputPicker outputpicker.Model
	ctx          *context.ProgramContext
}

func NewModel(ctx *context.ProgramContext) Model {
	lf := listfiles.NewModel(ctx)
	lf.SetTitle("Choose Order")

	outputFields := []outputpicker.Field{
		{
			Placeholder: "./merged.pdf",
			Prompt:      "Output File: ",
		},
	}
	op := outputpicker.NewModel(outputFields)
	return Model{
		fileList:     lf,
		outputPicker: op,
		ctx:          ctx,
	}
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
	// var fileListView, outputPickerView string
	var fileListStyle, outputPickerStyle lipgloss.Style
	if m.focusIndex == 0 {
		// fileListView = style.DefaultStyle.FocusedBorder.Render(m.fileList.View())
		// outputPickerView = style.DefaultStyle.BlurredBorder.Render(m.outputPicker.View())
		fileListStyle = style.DefaultStyle.FocusedBorder
		outputPickerStyle = style.DefaultStyle.BlurredBorder
	} else {
		// fileListView = style.DefaultStyle.BlurredBorder.Render(m.fileList.View())
		// outputPickerView = style.DefaultStyle.FocusedBorder.Render(m.outputPicker.View())
		fileListStyle = style.DefaultStyle.BlurredBorder
		outputPickerStyle = style.DefaultStyle.FocusedBorder
	}

	return style.RenderTwoFullRows(
		m.ctx.ScreenWidth,
		m.ctx.ScreenHeight,
		fileListStyle,
		outputPickerStyle,
		m.fileList.View(),
		m.outputPicker.View(),

	)
	// return "\n" + lipgloss.JoinVertical(
	// 	lipgloss.Left,
	// 	fileListView,
	// 	outputPickerView,
	// )

	// return style.RenderCol umnLayout(m.ctx.ScreenWidth, 10, fileListView, outputPickerView)
}
