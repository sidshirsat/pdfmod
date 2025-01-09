package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sidshirsat/pdfmod/internal/utils"
)

type FilePickerService struct {
	Prompter utils.Prompter
}

var _ FileHandler = &FilePickerService{}

// listFiles lists all the files in a directory and returns their names.
func (f *FilePickerService) ListFiles(dir string) ([]os.FileInfo, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fileInfos := make([]os.FileInfo, len(files))
	for i, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		fileInfos[i] = info
	}
	return fileInfos, nil
}

// SelectFile allows the user to select a file by entering the corresponding number.
func (f *FilePickerService) SelectFile(files []os.FileInfo) (string, error) {
	// Filter valid PDF files
	var validPDFs []os.FileInfo
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			validPDFs = append(validPDFs, file)
		}
	}

	// Display valid PDF files
	for i, file := range validPDFs {
		fmt.Printf("[%d] %s\n", i+1, file.Name())
	}

	// Check if no valid PDFs were found
	if len(validPDFs) == 0 {
		return "", fmt.Errorf(utils.Colorize("no valid PDF files found. Consider adding files..", utils.Red))
	}

	var selection int
	for {
		// Prompt for selection
		fmt.Print(utils.Colorize("Select a file number: ", utils.Blue))
		_, err := fmt.Scan(&selection)
		if err != nil {
			return "", fmt.Errorf(utils.Colorize("invalid selection: %w", utils.Red), err)
		}

		// Validate the selection
		if selection >= 1 && selection <= len(validPDFs) {
			selectedFile := validPDFs[selection-1]
			return selectedFile.Name(), nil
		} else {
			fmt.Println(utils.Colorize("Invalid selection. Please select a valid file number.", utils.Red))
		}
	}
}

// RenameFile renames a file with a new name.
func (f *FilePickerService) RenameFile(filePath, newName string) (string, error) {
	newPath := filepath.Join(filepath.Dir(filePath), newName+".pdf")
	err := os.Rename(filePath, newPath)
	if err != nil {
		return "", fmt.Errorf("failed to rename file: %w", err)
	}

	return newPath, nil
}
