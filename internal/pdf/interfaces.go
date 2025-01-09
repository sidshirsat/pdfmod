package pdf

// PDFMetadataHandler defines methods for PDF metadata operations.
type PDFMetadataHandler interface {
	UpdateMetadata(filePath, title, producer string) error
}
