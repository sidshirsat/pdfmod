package pdf

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

const maxRetries = 10 // Set a limit on the number of retries

// PDFService is a service to update PDF metadata.
type PDFService struct{}

// NewPDFService creates a new PDFService instance.
func NewPDFService() *PDFService {
	return &PDFService{}
}

var _ PDFMetadataHandler = &PDFService{}

// UpdatePDFMetadata updates the title and producer name in the PDF metadata, with retry on failure.
func (s *PDFService) UpdateMetadata(filePath, title, name string) error {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Attempt %d to update PDF metadata", attempt)

		// Read the original PDF file.
		pdfData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("could not read PDF file: %w", err)
		}

		// Replace the /Title and /Producer fields with the new title and name.
		pdfData = s.replaceMetadataField(pdfData, "Title", title)
		pdfData = s.replaceMetadataField(pdfData, "Producer", name)

		// Write the updated data back to the PDF file.
		err = os.WriteFile(filePath, pdfData, 0644)
		if err != nil {
			return fmt.Errorf("could not write updated PDF file: %w", err)
		}

		// Verify the update by reading back the file and checking for the updated metadata
		if s.verifyUpdate(filePath, title, name) {
			log.Println("PDF metadata updated successfully")
			return nil
		}

		log.Printf("Verification failed for attempt %d, retrying...", attempt)
		time.Sleep(1 * time.Second) // Optional delay between retries
	}

	return fmt.Errorf("failed to update PDF metadata after %d attempts", maxRetries)
}

// replaceMetadataField replaces the metadata field in the PDF file data.
func (s *PDFService) replaceMetadataField(pdfData []byte, fieldName, newValue string) []byte {
	fieldPrefix := fmt.Sprintf("/%s (", fieldName)
	start := bytes.Index(pdfData, []byte(fieldPrefix))
	if start == -1 {
		log.Printf("Field %s not found in PDF. Skipping replacement.", fieldName)
		// Field not found; return original data unchanged.
		return pdfData
	}
	start += len(fieldPrefix)

	// Find the end of the current field value (marked by ')')
	end := bytes.Index(pdfData[start:], []byte(")"))
	if end == -1 {
		log.Printf("End of field %s value not found. Skipping replacement.", fieldName)
		return pdfData
	}
	end += start

	// Replace the field value with the new value
	log.Printf("Replacing %s field value with: %s", fieldName, newValue)
	return bytes.Replace(pdfData, pdfData[start:end], []byte(newValue), 1)
}

// verifyUpdate reads the file and checks if the title and producer fields were updated correctly.
func (s *PDFService) verifyUpdate(filePath, expectedTitle, expectedProducer string) bool {
	// Read the file back to verify the update
	pdfData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read PDF file for verification: %v", err)
		return false
	}

	// Check if the updated metadata values are present
	titleUpdated := bytes.Contains(pdfData, []byte(fmt.Sprintf("/Title (%s)", expectedTitle)))
	producerUpdated := bytes.Contains(pdfData, []byte(fmt.Sprintf("/Producer (%s)", expectedProducer)))

	if titleUpdated && producerUpdated {
		log.Println("Verification succeeded: Title and Producer fields updated correctly.")
		return true
	}

	log.Println("Verification failed: Title or Producer fields not updated correctly.")
	return false
}
