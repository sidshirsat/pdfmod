package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sidshirsat/pdfmod/internal/file"
	"github.com/sidshirsat/pdfmod/internal/pdf"
	"github.com/sidshirsat/pdfmod/internal/utils"
)

// PDFManager handles user interactions and operations on the PDF file.
type PDFManager struct {
	FileHandler        file.FileHandler
	PDFMetadataHandler pdf.PDFMetadataHandler
	Prompter           Prompter
}

func NewPDFManager(fh file.FileHandler, pmh pdf.PDFMetadataHandler, prompter Prompter) *PDFManager {
	return &PDFManager{
		FileHandler:        fh,
		PDFMetadataHandler: pmh,
		Prompter:           prompter,
	}
}

var _ PDFManagerInterface = &PDFManager{}

func (pm *PDFManager) Execute() error {
	// Get the directory and list files.
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current directory: %w", err)
	}

	dir = filepath.Join(dir, "pdf_files")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %q does not exist. Please create it and add some PDF files", dir)
	}

	files, err := pm.FileHandler.ListFiles(dir)
	if err != nil {
		return err
	}

	// Select a file
	selectedFile, err := pm.FileHandler.SelectFile(files)
	if err != nil {
		return err
	}
	filePath := filepath.Join(dir, selectedFile)

	// Ask user choice
	fmt.Println("What would you like to do with the PDF:")
	fmt.Println("1. Rename the PDF")
	fmt.Println("2. Modify PDF metadata fields")
	choice := pm.Prompter.PromptUser("Enter the number of your choice: ")

	switch choice {
	case "1":
		newName := pm.Prompter.PromptUser("Enter the new name for the PDF (without extension): ")
		_, err = pm.FileHandler.RenameFile(filePath, newName)
		if err != nil {
			return err
		}
		fmt.Println(utils.Colorize("File renamed successfully.", utils.Green))
	case "2":
		title := pm.Prompter.PromptUser("Enter the new title for the PDF: ")
		producer := pm.Prompter.PromptUser("Enter the new producer name for the PDF: ")
		err := pm.PDFMetadataHandler.UpdateMetadata(filePath, title, producer)
		if err != nil {
			return err
		}
		fmt.Println(utils.Colorize("PDF metadata updated successfully.", utils.Green))
	default:
		fmt.Println(utils.Colorize("Invalid choice. Please restart and select '1' or '2'.", utils.Red))
		return fmt.Errorf("invalid choice: %s", choice) // Return an error for invalid choice
	}
	return nil
}
