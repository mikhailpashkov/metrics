package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type FileBackupRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewFileBackupRepository(filePath string) *FileBackupRepository {
	return &FileBackupRepository{filePath: filePath}
}

func (r *FileBackupRepository) FindAll(ctx context.Context) ([]*models.BackupMetrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]*models.BackupMetrics, 0)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if string(fileContent) == "" {
		fileContent = []byte("[]")
	}

	err = json.Unmarshal(fileContent, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file content: %w", err)
	}

	return result, nil
}

func (r *FileBackupRepository) SaveAll(ctx context.Context, metrics []*models.BackupMetrics) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sort.Slice(metrics, func(i, j int) bool { return metrics[i].ID > metrics[j].ID })

	tempFilePath := r.filePath + ".tmp"
	file, err := os.OpenFile(tempFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open temp file: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	err = enc.Encode(metrics)
	if err != nil {
		return fmt.Errorf("failed to encode metrics: %w", err)
	}

	err = os.Rename(tempFilePath, r.filePath)
	if err != nil {
		return fmt.Errorf("failed to replace target file with temp file: %w", err)
	}

	return nil
}
