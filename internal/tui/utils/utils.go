package utils

// TODO:
// show all errors when batch encrypt and decrypt using floating component
// show err when saving a file that already exists

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
		outFile = addExtension(outFile, ".pdf")
		err := api.MergeCreateFile(inFiles, getNextAvailablePath(outFile), false, nil)
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
				outFile = getNextAvailablePath(outFile)
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

func Decrypt(inFiles []string, password, outFilePath, outFilePrefix string, inPlace bool) tea.Cmd {
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
				outFile = getNextAvailablePath(outFile)
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

func Split(inFiles []string, outFilePath, outFilePrefix string, selectedPages []string, mergeIntoOne, extractAll bool) tea.Cmd {
	return func() tea.Msg {
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
			outFile = getNextAvailablePath(outFile)
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
			outFile = getNextAvailablePath(outFile)
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

func Img2Pdf(inFiles []string, outFile string, mergeIntoOne bool) tea.Cmd {
	return func() tea.Msg {
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
		
		outFile = addExtension(outFile, ".pdf")
		if mergeIntoOne {
			err := api.ImportImagesFile(inFiles, getNextAvailablePath(outFile), nil, nil)
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
			err := api.ImportImagesFile([]string{inFile}, getNextAvailablePath(outFile), nil, nil)
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

func ExtractImgs(inFiles []string, outFilePath string) tea.Cmd {
	return func() tea.Msg {
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

func getNextAvailablePath(basePath string) string {
	if _, err := os.Stat(basePath); err != nil {
		return basePath
	}

	ext := filepath.Ext(basePath)
	nameWithoutExt := strings.TrimSuffix(basePath, ext)

	counter := 1
	for {
		newPath := fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}

func addExtension(path, newExt string) string {
	ext := filepath.Ext(path)
	switch ext {
	case newExt:
		return path
	case "":
		return path + newExt
	}

	nameWithoutExt := strings.TrimSuffix(path, ext)
	return nameWithoutExt + newExt
}
