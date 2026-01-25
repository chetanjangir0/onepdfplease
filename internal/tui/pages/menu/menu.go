package menu

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	title string
	page  types.Page
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	tools  list.Model
	choice string
	ctx    *context.ProgramContext
}

func NewModel(ctx *context.ProgramContext) Model {
	items := []list.Item{
		item{title: "Merge PDFs", page: types.MergePage},
		item{title: "Split PDF", page: types.SplitPage},
		item{title: "Encrypt PDFs", page: types.EncryptPage},
		item{title: "Decrypt PDFs", page: types.DecryptPage},
		item{title: "Image(s) to pdf", page: types.Img2PdfPage},
		item{title: "Extract Images", page: types.ExtractImgsPage},
	}

	const defaultWidth = 20
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What tool do you want to use?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)
	l.KeyMap.Quit.SetEnabled(false)
	return Model{
		tools: l,
		ctx:   ctx,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.tools.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.tools.SelectedItem().(item)
			if ok {
				// log.Println("navigating to:" + string(i.title))
				return m, func() tea.Msg {
					return messages.Navigate{Page: i.page}
				}
			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.tools, cmd = m.tools.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	toolsView := style.DefaultStyle.FocusedBorder.Render(m.tools.View())
	return lipgloss.Place(
		m.ctx.TermWidth,
		m.ctx.MainContentHeight,
		lipgloss.Center,
		lipgloss.Center,
		toolsView,
	)
}
