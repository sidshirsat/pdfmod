package pdf_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/sidshirsat/pdfmod/internal/pdf"
)

func TestUpdateMetadata(t *testing.T) {
	// Create a temporary PDF file for testing
	tempFile, err := os.CreateTemp("", "test.pdf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some dummy PDF content to the temp file
	dummyPDFContent := []byte("%PDF-1.4\n1 0 obj\n<< /Title (Old Title) /Producer (Old Producer) >>\nendobj\ntrailer\n<< /Root 1 0 R >>\n%%EOF")
	if _, err := tempFile.Write(dummyPDFContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Create a new PDFService instance
	service := pdf.NewPDFService()

	// Define test cases
	tests := []struct {
		title    string
		producer string
	}{
		{"New Title", "New Producer"},
		{"Another Title", "Another Producer"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			// Call the UpdateMetadata method
			err := service.UpdateMetadata(tempFile.Name(), tt.title, tt.producer)
			if err != nil {
				t.Fatalf("UpdateMetadata failed: %v", err)
			}

			// Read the updated file
			updatedPDFContent, err := os.ReadFile(tempFile.Name())
			if err != nil {
				t.Fatalf("Failed to read updated file: %v", err)
			}

			// Check if the metadata was updated correctly
			if !bytes.Contains(updatedPDFContent, []byte("/Title ("+tt.title+")")) {
				t.Errorf("Title not updated correctly, got: %s", updatedPDFContent)
			}
			if !bytes.Contains(updatedPDFContent, []byte("/Producer ("+tt.producer+")")) {
				t.Errorf("Producer not updated correctly, got: %s", updatedPDFContent)
			}
		})
	}
}
