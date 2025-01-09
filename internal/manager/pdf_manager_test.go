package manager_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sidshirsat/pdfmod/internal/manager"
	"github.com/sidshirsat/pdfmod/mocks"
)

func TestPDFManager_Execute_RenameFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Resolve the absolute path to "pdf_files" in the project root
	absPDFDir, err := filepath.Abs("pdf_files")
	if err != nil {
		t.Fatalf("failed to resolve absolute path for pdf_files: %v", err)
	}

	// Ensure the directory exists
	err = os.MkdirAll(absPDFDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create pdf_files directory: %v", err)
	}

	mockFileHandler := mocks.NewMockFileHandler(ctrl)
	mockPDFMetadataHandler := mocks.NewMockPDFMetadataHandler(ctrl)
	mockPrompter := mocks.NewMockPrompter(ctrl)
	mockFileInfo := mocks.NewMockFileInfo(ctrl)

	// Configure mockFileInfo to simulate a file with a name
	mockFileInfo.EXPECT().Name().Return("sample.pdf").AnyTimes()
	mockFileInfo.EXPECT().IsDir().Return(false).AnyTimes()

	// Set expectations for FileHandler and Prompter interactions
	mockFileHandler.EXPECT().ListFiles(absPDFDir).Return([]os.FileInfo{mockFileInfo}, nil).Times(1)
	mockFileHandler.EXPECT().SelectFile([]os.FileInfo{mockFileInfo}).Return("sample.pdf", nil).Times(1)

	mockPrompter.EXPECT().PromptUser("Enter the number of your choice: ").Return("1").Times(1)
	mockPrompter.EXPECT().PromptUser("Enter the new name for the PDF (without extension): ").Return("new_sample").Times(1)
	mockFileHandler.EXPECT().RenameFile(filepath.Join(absPDFDir, "sample.pdf"), "new_sample").Return("new_sample.pdf", nil).Times(1)

	// Create PDFManager instance with mocks
	pdfManager := manager.NewPDFManager(mockFileHandler, mockPDFMetadataHandler, mockPrompter)

	// Execute and assert no error
	if err := pdfManager.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Clean up
	_ = os.RemoveAll(absPDFDir)
}

func TestPDFManager_Execute_UpdateMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	absPDFDir, err := filepath.Abs("pdf_files")
	if err != nil {
		t.Fatalf("failed to resolve absolute path for pdf_files: %v", err)
	}

	err = os.MkdirAll(absPDFDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create pdf_files directory: %v", err)
	}

	mockFileHandler := mocks.NewMockFileHandler(ctrl)
	mockPDFMetadataHandler := mocks.NewMockPDFMetadataHandler(ctrl)
	mockPrompter := mocks.NewMockPrompter(ctrl)
	mockFileInfo := mocks.NewMockFileInfo(ctrl)

	mockFileInfo.EXPECT().Name().Return("sample.pdf").AnyTimes()
	mockFileInfo.EXPECT().IsDir().Return(false).AnyTimes()

	mockFileHandler.EXPECT().ListFiles(absPDFDir).Return([]os.FileInfo{mockFileInfo}, nil).Times(1)
	mockFileHandler.EXPECT().SelectFile([]os.FileInfo{mockFileInfo}).Return("sample.pdf", nil).Times(1)

	mockPrompter.EXPECT().PromptUser("Enter the number of your choice: ").Return("2").Times(1)
	mockPrompter.EXPECT().PromptUser("Enter the new title for the PDF: ").Return("New Title").Times(1)
	mockPrompter.EXPECT().PromptUser("Enter the new producer name for the PDF: ").Return("New Producer").Times(1)

	mockPDFMetadataHandler.EXPECT().UpdateMetadata(filepath.Join(absPDFDir, "sample.pdf"), "New Title", "New Producer").Return(nil).Times(1)

	pdfManager := manager.NewPDFManager(mockFileHandler, mockPDFMetadataHandler, mockPrompter)

	if err := pdfManager.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_ = os.RemoveAll(absPDFDir)
}

func TestPDFManager_Execute_InvalidChoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Resolve the absolute path to "pdf_files" in the project root
	absPDFDir, err := filepath.Abs("pdf_files")
	if err != nil {
		t.Fatalf("failed to resolve absolute path for pdf_files: %v", err)
	}

	// Ensure the directory exists
	err = os.MkdirAll(absPDFDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create pdf_files directory: %v", err)
	}

	mockFileHandler := mocks.NewMockFileHandler(ctrl)
	mockPDFMetadataHandler := mocks.NewMockPDFMetadataHandler(ctrl)
	mockPrompter := mocks.NewMockPrompter(ctrl)
	mockFileInfo := mocks.NewMockFileInfo(ctrl)

	// Configure mockFileInfo to simulate a file with a name
	mockFileInfo.EXPECT().Name().Return("sample.pdf").AnyTimes()
	mockFileInfo.EXPECT().IsDir().Return(false).AnyTimes()

	// Set expectations for FileHandler and Prompter interactions
	mockFileHandler.EXPECT().ListFiles(absPDFDir).Return([]os.FileInfo{mockFileInfo}, nil).Times(1)
	mockFileHandler.EXPECT().SelectFile([]os.FileInfo{mockFileInfo}).Return("sample.pdf", nil).Times(1)

	mockPrompter.EXPECT().PromptUser("Enter the number of your choice: ").Return("3").Times(1)

	// Create PDFManager instance with mocks
	pdfManager := manager.NewPDFManager(mockFileHandler, mockPDFMetadataHandler, mockPrompter)

	// Execute and assert error
	if err := pdfManager.Execute(); err == nil {
		t.Fatalf("expected an error, got none")
	}

	// Clean up
	_ = os.RemoveAll(absPDFDir)
}
