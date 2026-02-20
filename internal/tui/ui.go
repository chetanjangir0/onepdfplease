package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chetanjangir0/onepdfplease/internal/tui/components/footer"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/decrypt"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/doc2pdf"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/encrypt"
	extractimgs "github.com/chetanjangir0/onepdfplease/internal/tui/pages/extractImgs"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/img2pdf"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/menu"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/merge"
	"github.com/chetanjangir0/onepdfplease/internal/tui/pages/split"
	"github.com/chetanjangir0/onepdfplease/internal/tui/style"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

type model struct {
	quitting bool
	ctx      *context.ProgramContext

	// each page has its own model
	menuModel        menu.Model
	mergeModel       merge.Model
	encryptModel     encrypt.Model
	decryptModel     decrypt.Model
	splitModel       split.Model
	img2pdfModel     img2pdf.Model
	extractImgsModel extractimgs.Model
	doc2pdfModel     doc2pdf.Model
	footer           footer.Model
}

func InitialModel() model {
	m := model{}
	m.ctx = &context.ProgramContext{
		StatusType:  context.None,
		CurrentPage: types.MenuPage,
	}
	m.menuModel = menu.NewModel(m.ctx)
	m.mergeModel = merge.NewModel(m.ctx)
	m.encryptModel = encrypt.NewModel(m.ctx)
	m.decryptModel = decrypt.NewModel(m.ctx)
	m.splitModel = split.NewModel(m.ctx)
	m.img2pdfModel = img2pdf.NewModel(m.ctx)
	m.extractImgsModel = extractimgs.NewModel(m.ctx)
	m.doc2pdfModel = doc2pdf.NewModel(m.ctx)
	m.footer = footer.NewModel(m.ctx)
	return m
}

func (m model) Init() tea.Cmd {
	return m.footer.Init()
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)
	case tea.KeyMsg:
		if m.ctx.StatusType != context.Processing {
			m.ctx.ClearStatus()
		}

		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "?":
			if !m.footer.ShowAll {
				m.ctx.MainContentHeight = m.ctx.MainContentHeight +
					style.FooterHeight - style.ExpandedHelpHeight
			} else {
				m.ctx.MainContentHeight = m.ctx.MainContentHeight +
					style.ExpandedHelpHeight - style.FooterHeight
			}
			m.footer.ShowAll = !m.footer.ShowAll
		}
	case messages.Navigate:
		m.ctx.CurrentPage = msg.Page
		return m, nil
	case messages.PDFOperationStatus:
		if msg.Err != nil {
			m.ctx.SetStatusError(fmt.Sprintf("Failed: %v", msg.Err))
		} else {
			m.ctx.SetStatusSuccess(fmt.Sprintf("%s completed successfully", msg.TaskType))
		}
	case messages.ShowError:
		if msg.Err != nil {
			m.ctx.SetStatusError(fmt.Sprintf("Error: %v", msg.Err))
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	switch m.ctx.CurrentPage {
	case types.MenuPage:
		m.menuModel, cmd = m.menuModel.Update(msg)
	case types.MergePage:
		m.mergeModel, cmd = m.mergeModel.Update(msg)
	case types.EncryptPage:
		m.encryptModel, cmd = m.encryptModel.Update(msg)
	case types.DecryptPage:
		m.decryptModel, cmd = m.decryptModel.Update(msg)
	case types.SplitPage:
		m.splitModel, cmd = m.splitModel.Update(msg)
	case types.Img2PdfPage:
		m.img2pdfModel, cmd = m.img2pdfModel.Update(msg)
	case types.ExtractImgsPage:
		m.extractImgsModel, cmd = m.extractImgsModel.Update(msg)
	case types.Doc2PdfPage:
		m.doc2pdfModel, cmd = m.doc2pdfModel.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var view string
	switch m.ctx.CurrentPage {
	case types.MenuPage:
		view = m.menuModel.View()
	case types.MergePage:
		view = m.mergeModel.View()
	case types.EncryptPage:
		view = m.encryptModel.View()
	case types.DecryptPage:
		view = m.decryptModel.View()
	case types.SplitPage:
		view = m.splitModel.View()
	case types.Img2PdfPage:
		view = m.img2pdfModel.View()
	case types.ExtractImgsPage:
		view = m.extractImgsModel.View()
	case types.Doc2PdfPage:
		view = m.doc2pdfModel.View()
	}
	return lipgloss.JoinVertical(lipgloss.Left, view, m.footer.View())

}

func (m *model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.ctx.TermWidth = msg.Width
	m.ctx.TermHeight = msg.Height
	m.ctx.MainContentHeight = msg.Height
	if m.footer.ShowAll {
		m.ctx.MainContentHeight = msg.Height - style.ExpandedHelpHeight
	} else {
		m.ctx.MainContentHeight = msg.Height - style.FooterHeight
	}
}
