package utils

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func Merge(inFiles []string, outFile string) tea.Cmd {
	return func() tea.Msg {
		tasktype := "Merge"
		if len(inFiles) == 0 {
			return messages.PDFOperationStatus{
				TaskType: tasktype,
				Err:      fmt.Errorf("There are no files to merge"),
			}
		}
		err := api.MergeCreateFile(inFiles, outFile, false, nil)
		if err != nil {
			return messages.PDFOperationStatus{
				TaskType: tasktype,
				Err:      err,
			}
		}

		return messages.PDFOperationStatus{
			TaskType: tasktype,
		}
	}
}
