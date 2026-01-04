package merge

// TODO:
// change focus using tab
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/listfiles"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/outputpicker"
)

type Model struct {
	focusIndex   int // 0 for fileList 1 for outputPicker
	fileList     listfiles.Model
	outputPicker outputpicker.Model
}

func NewModel() Model {
	lf := listfiles.NewModel()
	lf.Title = "Merge PDFs"

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
	outputPickerView := m.outputPicker.View()
	fileListView := m.fileList.View()

	return "\n" + lipgloss.JoinVertical(
		lipgloss.Left,
		fileListView,
		outputPickerView,
	)
}
