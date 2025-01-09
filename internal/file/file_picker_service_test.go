package file_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sidshirsat/pdfmod/internal/file"
)

func TestFilePickerService_ListFiles_DirectoryNotFound(t *testing.T) {
	fps := &file.FilePickerService{}

	// Attempt to list files in a non-existent directory
	_, err := fps.ListFiles("non_existent_dir")
	if err == nil {
		t.Fatal("Expected error for non-existent directory, got nil")
	}
}

func TestFilePickerService_SelectFile_Success(t *testing.T) {
	// Resolve the absolute path to "pdf_files" in the project root
	absPDFDir, err := filepath.Abs("pdf_files")
	if err != nil {
		t.Fatalf("Failed to resolve absolute path for pdf_files: %v", err)
	}

	// Ensure the directory exists
	err = os.MkdirAll(absPDFDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create pdf_files directory: %v", err)
	}
	defer os.RemoveAll(absPDFDir) // Clean up after test

	// Create mock files in the pdf_files directory
	fileNames := []string{"file1.pdf", "file2.pdf", "file3.txt"}
	for _, name := range fileNames {
		_, err := os.Create(filepath.Join(absPDFDir, name))
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", name, err)
		}
	}

	// Create FilePickerService and call ListFiles
	fps := &file.FilePickerService{}
	files, err := fps.ListFiles(absPDFDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Call SelectFile
	selectedFile, err := fps.SelectFile(files)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the correct file was selected
	if selectedFile != "file1.pdf" {
		t.Errorf("Expected file1.pdf, got %s", selectedFile)
	}
}

func TestFilePickerService_RenameFile_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a temporary file within the temp directory
	tempFile, err := os.CreateTemp(tempDir, "testfile-*.pdf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()

	// Define the new name for the file
	newName := "renamed-file"

	// Initialize FilePickerService
	fps := &file.FilePickerService{}

	// Call RenameFile
	newFilePath, err := fps.RenameFile(tempFilePath, newName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the file was renamed correctly
	expectedPath := filepath.Join(tempDir, newName+".pdf")
	if newFilePath != expectedPath {
		t.Errorf("Expected new file path to be '%s', got '%s'", expectedPath, newFilePath)
	}

	// Check that the new file exists
	if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
		t.Errorf("Expected renamed file to exist at '%s', but it does not", newFilePath)
	}

	// Check that the old file no longer exists
	if _, err := os.Stat(tempFilePath); !os.IsNotExist(err) {
		t.Errorf("Expected original file to be moved, but it still exists at '%s'", tempFilePath)
	}
}

func TestFilePickerService_RenameFile_Failure(t *testing.T) {
	// Initialize FilePickerService
	fps := &file.FilePickerService{}

	// Attempt to rename a non-existent file
	nonExistentFile := "nonexistent.pdf"
	newName := "should-not-exist"

	_, err := fps.RenameFile(nonExistentFile, newName)
	if err == nil {
		t.Fatal("Expected error when renaming non-existent file, got nil")
	}

	// Check if the underlying error is "file does not exist"
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected a file does not exist error, got %v", err)
	}
}
