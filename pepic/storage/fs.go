package storage

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

var errOutsideStorage = errors.New("path escapes the storage directory")

type FileSystemBackend struct {
	dir string
}

func NewFileSystemBackend(dir string) *FileSystemBackend {
	fs := new(FileSystemBackend)
	fs.dir = dir
	return fs
}

// fullPath resolves objectName inside fs.dir and refuses anything that would
// land outside it, whatever the caller passed in.
func (fs *FileSystemBackend) fullPath(objectName string) (string, error) {
	root := filepath.Clean(fs.dir)
	full := filepath.Join(root, objectName)
	rel, err := filepath.Rel(root, full)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errOutsideStorage
	}
	return full, nil
}

func (fs *FileSystemBackend) PutObject(objectName string, data []byte) (string, error) {
	fullPath, err := fs.fullPath(objectName)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return "", err
	}

	dst, err := os.Create(fullPath)

	if err != nil {
		return "", err
	}

	defer dst.Close()

	if _, err = dst.Write(data); err != nil {
		return "", err
	}

	return fullPath, nil
}

func (fs *FileSystemBackend) GetObject(objectName string) ([]byte, error) {
	fullPath, err := fs.fullPath(objectName)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fs *FileSystemBackend) IsExists(objectName string) bool {
	fullPath, err := fs.fullPath(objectName)
	if err != nil {
		return false
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func (fs *FileSystemBackend) Size(objectName string) int64 {
	fullPath, err := fs.fullPath(objectName)
	if err != nil {
		return 0
	}
	file, err := os.Open(fullPath)
	if err != nil {
		return 0
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return 0
	}
	return info.Size()
}

func (fs *FileSystemBackend) Proxy(c echo.Context, objectName string) error {
	fullPath, err := fs.fullPath(objectName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "File not found")
	}
	return c.File(fullPath)
}
