package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BaseDir string
}

func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{BaseDir: baseDir}, nil
}

func (s *LocalStorage) UploadFile(ctx context.Context, sessionID, fileName string, content io.Reader) (string, error) {
	// Generate unique path
	// Structure: {session_id}/YYYY/MM/DD/uuid-filename
	if sessionID == "" {
		sessionID = "default"
	}
	now := time.Now()
	subDir := filepath.Join(sessionID, now.Format("2006"), now.Format("01"), now.Format("02"))
	fullDir := filepath.Join(s.BaseDir, subDir)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", err
	}

	uniqueName := uuid.New().String() + "-" + fileName
	fullPath := filepath.Join(fullDir, uniqueName)
	relPath := filepath.Join(subDir, uniqueName)

	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, content)
	if err != nil {
		return "", err
	}

	return relPath, nil
}

func (s *LocalStorage) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.BaseDir, path)
	return os.Open(fullPath)
}
func (s *LocalStorage) DeleteSessionDir(ctx context.Context, sessionID string) error {
	if sessionID == "" || sessionID == "default" {
		return nil // Safety: don't delete default or empty session dir easily via this
	}
	path := filepath.Join(s.BaseDir, sessionID)
	return os.RemoveAll(path)
}
