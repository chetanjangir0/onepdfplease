package utils

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func Merge(inFiles []string, outFile string) tea.Cmd{
	return func() tea.Msg {
		err := api.MergeCreateFile(inFiles, outFile, false, nil)
		if err != nil {
			return messages.PDFOperationStatus{
				TaskType: "PDF merging",
				Err: err,
			}
		}

		return messages.PDFOperationStatus{
			TaskType: "PDFmerging",
		}
	}
}
