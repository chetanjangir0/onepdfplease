package pdf

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func Merge(inFiles []string, outFile string, ctx *context.ProgramContext) tea.Cmd {
	return func() tea.Msg {
		ctx.SetStatusProcessing("merging files...")
		taskType := "Merge"

		for _, f := range inFiles {
			if _, err := os.Stat(f); err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      fmt.Errorf("file not found: %s", f),
				}
			}
		}

		if len(inFiles) <= 1 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("At least 2 files required for merge"),
			}
		}
		outFile = utils.UpdateFileExtension(outFile, ".pdf")
		err := api.MergeCreateFile(inFiles, utils.GetNextAvailablePath(outFile), false, nil)
		if err != nil {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      err,
			}
		}

		return messages.PDFOperationStatus{
			TaskType: taskType,
		}
	}
}
