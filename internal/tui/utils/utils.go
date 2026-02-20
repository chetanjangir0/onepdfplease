package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
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
		outFile = updateFileExtension(outFile, ".pdf")
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

		outFile = updateFileExtension(outFile, ".pdf")
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

func Doc2Pdf(inFiles []string, outDir, outPrefix string, ctx *context.ProgramContext) tea.Cmd {
	return func() tea.Msg {
		ctx.SetStatusProcessing("converting documents...")
		taskType := "Doc to PDF"
		successMsg := messages.PDFOperationStatus{TaskType: taskType}

		if len(inFiles) == 0 {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("There are no documents to convert"),
			}
		}

		bin, err := findLibreOfficeBinary()
		if err != nil {
			return messages.PDFOperationStatus{TaskType: taskType, Err: err}
		}

		if outDir == "" {
			outDir = "./"
		}
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("Invalid output path: %v", err),
			}
		}
		if st, err := os.Stat(outDir); err != nil || !st.IsDir() {
			return messages.PDFOperationStatus{
				TaskType: taskType,
				Err:      fmt.Errorf("Output path is not a directory: %s", outDir),
			}
		}

		for _, f := range inFiles {
			if _, err := os.Stat(f); err != nil {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      fmt.Errorf("file not found: %s", f),
				}
			}
			ext := strings.ToLower(filepath.Ext(f))
			if ext != ".doc" && ext != ".docx" {
				return messages.PDFOperationStatus{
					TaskType: taskType,
					Err:      fmt.Errorf("unsupported file type: %s", filepath.Base(f)),
				}
			}
		}

		var failedFiles []string
		for _, inFile := range inFiles {
			tmpDir, err := os.MkdirTemp("", "onepdfplease-doc2pdf-*")
			if err != nil {
				failedFiles = append(failedFiles, filepath.Base(inFile))
				continue
			}

			profileDir := filepath.Join(tmpDir, "lo-profile")
			_ = os.MkdirAll(profileDir, 0755)
			profileURI := "file://" + filepath.ToSlash(profileDir)

			cmd := exec.Command(
				bin,
				"-env:UserInstallation="+profileURI,
				"--headless",
				"--nologo",
				"--nolockcheck",
				"--nodefault",
				"--norestore",
				"--convert-to",
				"pdf",
				"--outdir",
				tmpDir,
				inFile,
			)

			out, err := cmd.CombinedOutput()
			if err != nil {
				_ = os.RemoveAll(tmpDir)
				errMsg := strings.TrimSpace(string(out))
				if len(errMsg) > 0 {
					_ = errMsg
				}
				failedFiles = append(failedFiles, filepath.Base(inFile))
				continue
			}

			base := strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile))
			src := filepath.Join(tmpDir, base+".pdf")
			if _, err := os.Stat(src); err != nil {
				_ = os.RemoveAll(tmpDir)
				failedFiles = append(failedFiles, filepath.Base(inFile))
				continue
			}

			dst := filepath.Join(outDir, outPrefix+base+".pdf")
			dst = getNextAvailablePath(dst)
			if err := moveFile(src, dst); err != nil {
				_ = os.RemoveAll(tmpDir)
				failedFiles = append(failedFiles, filepath.Base(inFile))
				continue
			}

			_ = os.RemoveAll(tmpDir)
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

func findLibreOfficeBinary() (string, error) {
	candidates := []string{"soffice", "libreoffice"}
	for _, c := range candidates {
		if p, err := exec.LookPath(c); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("LibreOffice not found. Please install it and ensure 'soffice' (or 'libreoffice') is on PATH")
}

func moveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	if err := copyFile(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
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

		if _, err := os.Stat(newPath); errors.Is(err, fs.ErrNotExist) {
			return newPath
		}
		counter++
	}
}

func updateFileExtension(path, newExt string) string {
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
