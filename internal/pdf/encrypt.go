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
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func Encrypt(
	inFiles []string,
	password,
	outFilePath,
	outFilePrefix string,
	inPlace bool,
	ctx *context.ProgramContext,
) tea.Cmd {
	return func() tea.Msg {
		ctx.SetStatusProcessing("encrypting files...")
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

		for _, f := range inFiles {
			if _, err := os.Stat(f); err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      fmt.Errorf("file not found: %s", f),
				}
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
				outFile = utils.GetNextAvailablePath(outFile)
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

func Decrypt(
	inFiles []string,
	password,
	outFilePath,
	outFilePrefix string,
	inPlace bool,
	ctx *context.ProgramContext,
) tea.Cmd {
	return func() tea.Msg {
		ctx.SetStatusProcessing("decrypting files...")
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

		for _, f := range inFiles {
			if _, err := os.Stat(f); err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      fmt.Errorf("file not found: %s", f),
				}
			}
		}

		conf := model.NewDefaultConfiguration()
		conf.UserPW = password

		var failedFiles []string
		for _, f := range inFiles {
			var outFile string
			if !inPlace {
				outFile = filepath.Join(outFilePath, outFilePrefix+filepath.Base(f))
				outFile = utils.GetNextAvailablePath(outFile)
			} else {
				outFile = ""
			}
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
