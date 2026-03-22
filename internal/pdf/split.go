package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func Split(
	inFiles []string,
	outFilePath,
	outFilePrefix string,
	selectedPages []string,
	mergeIntoOne,
	extractAll bool,
	ctx *context.ProgramContext,
) tea.Cmd {
	return func() tea.Msg {
		ctx.SetStatusProcessing("splitting file...")
		taskType := "Splitting"
		successMsg := messages.PDFOperationStatus{
			TaskType: taskType,
		}

		if len(inFiles) > 1 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("You can only split one file at a time"),
			}
		} else if len(inFiles) == 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("Please add a file first"),
			}
		}
		inFile := inFiles[0]

		if _, err := os.Stat(inFile); err != nil {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("file not found: %s", inFile),
			}
		}

		if extractAll {
			if err := api.SplitFile(inFile, outFilePath, 1, nil); err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      err,
				}
			}

			return successMsg
		}

		if mergeIntoOne {
			outFile := filepath.Join(outFilePath, outFilePrefix+filepath.Base(inFile))
			outFile = utils.GetNextAvailablePath(outFile)
			if err := api.TrimFile(inFile, outFile, selectedPages, nil); err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      err,
				}
			}

			return successMsg
		}

		var failedRanges []string
		for _, pageRange := range selectedPages {
			outFile := filepath.Join(outFilePath, outFilePrefix+pageRange+".pdf")
			outFile = utils.GetNextAvailablePath(outFile)
			if err := api.TrimFile(inFile, outFile, []string{pageRange}, nil); err != nil {
				failedRanges = append(failedRanges, pageRange)
			}
		}

		if len(failedRanges) > 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err: fmt.Errorf(
					"Failed to split %d range(s): %s",
					len(failedRanges),
					strings.Join(failedRanges, ","),
				),
			}
		}

		return successMsg
	}
}
