package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

// File file service
type File struct {
	mutex      sync.RWMutex
	fileFolder string
	files      map[string]*FileInfo
}

// FileInfo info about file
type FileInfo struct {
	Type string
	Path string
}

// NewFile new file service
func NewFile(fileFolder string) *File {
	return &File{
		fileFolder: fileFolder,
		files:      make(map[string]*FileInfo),
	}
}

// StoreFile service store file
func (s *File) StoreFile(fileType string, fileData bytes.Buffer) (string, error) {
	fileID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate file id: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s%s", s.fileFolder, fileID, fileType)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot create file file: %w", err)
	}

	_, err = fileData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write file to file: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.files[fileID.String()] = &FileInfo{
		Type: fileType,
		Path: filePath,
	}

	return fileID.String(), nil
}
