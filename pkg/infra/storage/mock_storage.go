package storage

import (
	"context"
	"errors"
)

type mockStorageAgent struct {
	repo map[string][]byte
}

func newMockStorageAgent() StorageAgent {
	return &mockStorageAgent{
		repo: make(map[string][]byte),
	}
}

func (s *mockStorageAgent) UploadFile(ctx context.Context, file []byte, filename string) error {
	s.repo[filename] = file
	return nil
}

func (s *mockStorageAgent) DownloadFile(ctx context.Context, filename string) ([]byte, error) {
	if filename == "testfilename" {
		return []byte("this is the test file"), nil
	}

	file, exist := s.repo[filename]
	if !exist {
		return nil, errors.New("file not found")
	}

	return file, nil
}
