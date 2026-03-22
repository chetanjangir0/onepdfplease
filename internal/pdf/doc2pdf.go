package pdf

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui/context"
	"github.com/chetanjangir0/onepdfplease/internal/tui/messages"
	"github.com/chetanjangir0/onepdfplease/internal/tui/utils"
)

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
			dst = utils.GetNextAvailablePath(dst)
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

