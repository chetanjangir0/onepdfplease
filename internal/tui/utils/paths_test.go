package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateFileExtension(t *testing.T) {
	tests := []struct {
		name         string
		path, newExt string
		expected     string
	}{
		{"No path", "", ".pdf", ".pdf"},
		{"No newExt", "/home/user/test.txt", "", "/home/user/test"},
		{"Normal", "/home/test.txt", ".pdf", "/home/test.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpdateFileExtension(tt.path, tt.newExt)
			if result != tt.expected {
				t.Errorf("Got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestGetNextAvailablePath(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		basePath      string
		existingFiles []string
		expected      string
	}{
		{
			name:          "file doesn't exists",
			basePath:      filepath.Join(tmpDir, "test.pdf"),
			existingFiles: []string{},
			expected:      filepath.Join(tmpDir, "test.pdf"),
		},
		{
			name:          "file exists, no collision",
			basePath:      filepath.Join(tmpDir, "test2.pdf"),
			existingFiles: []string{"test2.pdf"},
			expected:      filepath.Join(tmpDir, "test2_1.pdf"),
		},
		{
			name:          "file exists, multiple collisions",
			basePath:      filepath.Join(tmpDir, "test3.pdf"),
			existingFiles: []string{"test3.pdf", "test3_1.pdf", "test3_2.pdf"},
			expected:      filepath.Join(tmpDir, "test3_3.pdf"),
		},
		{
			name:          "No extension",
			basePath:      filepath.Join(tmpDir, "test4"),
			existingFiles: []string{"test4"},
			expected:      filepath.Join(tmpDir, "test4_1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, filename := range tt.existingFiles {
				path := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			result := GetNextAvailablePath(tt.basePath)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}

			// cleanup
			for _, filename := range tt.existingFiles{
				os.Remove(filepath.Join(tmpDir, filename))
			}
		})
	}
}
