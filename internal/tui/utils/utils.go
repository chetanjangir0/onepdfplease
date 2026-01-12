package utils

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
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

func Encrypt(inFiles []string, password, outFilePath, outFilePrefix string) tea.Cmd {
	return func() tea.Msg {
		taskType := "Encrypt"
		if len(password) == 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("Please provide a password"),
			}

		}
		if len(inFiles) == 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("There are no files to encrypt"),
			}
		}
		conf := model.NewDefaultConfiguration()
		conf.UserPW = password
		conf.OwnerPW = password
		conf.EncryptUsingAES = true
		conf.EncryptKeyLength = 256

		for _, f := range inFiles {
			outFile := outFilePath + outFilePrefix + filepath.Base(f)
			err := api.EncryptFile(f, outFile, conf)
			if err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      err,
				}
			}
		}

		return messages.PDFOperationStatus{
			TaskType: taskType,
		}

	}

}
