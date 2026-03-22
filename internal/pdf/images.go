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
func Img2Pdf(inFiles []string, outFile string, mergeIntoOne bool, ctx *context.ProgramContext) tea.Cmd {
	return func() tea.Msg {
		ctx.SetStatusProcessing("converting files...")
		taskType := "Image to Pdf"
		successMsg := messages.PDFOperationStatus{
			TaskType: taskType,
		}

		for _, f := range inFiles {
			if _, err := os.Stat(f); err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      fmt.Errorf("file not found: %s", f),
				}
			}
		}

		if len(inFiles) == 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("There are no images to convert"),
			}
		}

		outFile = utils.UpdateFileExtension(outFile, ".pdf")
		if mergeIntoOne {
			err := api.ImportImagesFile(inFiles, utils.GetNextAvailablePath(outFile), nil, nil)
			if err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      err,
				}
			}
			return successMsg
		}

		var failedFiles []string
		for _, inFile := range inFiles {
			err := api.ImportImagesFile([]string{inFile}, utils.GetNextAvailablePath(outFile), nil, nil)
			if err != nil {
				failedFiles = append(failedFiles, filepath.Base(inFile))
			}
		}

		if len(failedFiles) > 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err: fmt.Errorf(
					"Failed to convert %d file(s): %s",
					len(failedFiles),
					strings.Join(failedFiles, ","),
				),
			}
		}

		return successMsg

	}
}

func ExtractImgs(inFiles []string, outFilePath string, ctx *context.ProgramContext) tea.Cmd {
	return func() tea.Msg {
		ctx.SetStatusProcessing("extracting images...")
		taskType := "Extracting Images"
		successMsg := messages.PDFOperationStatus{
			TaskType: taskType,
		}

		if len(inFiles) > 1 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("You can only extract from one file at a time"),
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

		if err := api.ExtractImagesFile(inFile, outFilePath, nil, nil); err != nil {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      err,
			}
		}

		return successMsg

	}
}
