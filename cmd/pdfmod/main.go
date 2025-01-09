package main

import (
	"log"

	"github.com/sidshirsat/pdfmod/internal/file"
	"github.com/sidshirsat/pdfmod/internal/manager"
	"github.com/sidshirsat/pdfmod/internal/pdf"
	"github.com/sidshirsat/pdfmod/internal/utils"
)

func main() {
	// Initialize services

	pdfMetadataHandler := pdf.NewPDFService()
	basePrompter := &utils.BasePrompter{}
	prompter := &utils.ConsolePrompter{
		Prompter: basePrompter,
	}
	fileHandler := &file.FilePickerService{
		Prompter: prompter,
	}
	// Initialize PDF Manager
	pdfManager := manager.NewPDFManager(fileHandler, pdfMetadataHandler, prompter)

	// Execute the manager operation
	err := pdfManager.Execute()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
