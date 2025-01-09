package file

import "os"

// FileHandler defines methods for file operations.
type FileHandler interface {
	ListFiles(dir string) ([]os.FileInfo, error)
	SelectFile(files []os.FileInfo) (string, error)
	RenameFile(filePath, newName string) (string, error)
}
