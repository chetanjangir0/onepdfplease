package utils

// TODO:
// show all errors when batch encrypt and decrypt using floating component

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func Merge(inFiles []string, outFile string) tea.Cmd {
	return func() tea.Msg {
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
		err := api.MergeCreateFile(inFiles, outFile, false, nil)
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

func Encrypt(inFiles []string, password, outFilePath, outFilePrefix string, inPlace bool) tea.Cmd {
	return func() tea.Msg {
		taskType := "Encryption"
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

		var failedFiles []string
		for _, f := range inFiles {
			var outFile string
			if !inPlace {
				outFile = filepath.Join(outFilePath, outFilePrefix+filepath.Base(f))
			} else {
				outFile = ""
			}
			if err := api.EncryptFile(f, outFile, conf); err != nil {
				failedFiles = append(failedFiles, filepath.Base(f))
			}
		}
		if len(failedFiles) > 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err: fmt.Errorf(
					"Failed to encrypt %d file(s): %s",
					len(failedFiles),
					strings.Join(failedFiles, ","),
				),
			}
		}

		return messages.PDFOperationStatus{
			TaskType: taskType,
		}

	}

}

func Decrypt(inFiles []string, password, outFilePath, outFilePrefix string) tea.Cmd {
	return func() tea.Msg {
		taskType := "Decryption"
		if len(password) == 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("Please provide a password"),
			}

		}
		if len(inFiles) == 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("There are no files to decrypt"),
			}
		}
		conf := model.NewDefaultConfiguration()
		conf.UserPW = password

		var failedFiles []string
		for _, f := range inFiles {
			outFile := filepath.Join(outFilePath, outFilePrefix+filepath.Base(f))
			if err := api.DecryptFile(f, outFile, conf); err != nil {
				failedFiles = append(failedFiles, filepath.Base(f))
			}
		}
		if len(failedFiles) > 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err: fmt.Errorf(
					"Failed to decrypt %d file(s): %s",
					len(failedFiles),
					strings.Join(failedFiles, ","),
				),
			}
		}

		return messages.PDFOperationStatus{
			TaskType: taskType,
		}

	}

}

func Split(inFile, outFilePath, outFilePrefix string, selectedPages []string, mergeIntoOne bool) tea.Cmd {
	return func() tea.Msg {
		taskType := "Splitting"

		if mergeIntoOne {
			outFile := filepath.Join(outFilePath, outFilePrefix+filepath.Base(inFile))
			if err := api.TrimFile(inFile, outFile, selectedPages, nil); err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      err,
				}
			}
		} else {
			var failedRanges []string
			for _, pageRange := range selectedPages {
				outFile := filepath.Join(outFilePath, outFilePrefix+pageRange)
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
		}

		return messages.PDFOperationStatus{
			TaskType: taskType,
		}

	}
}
