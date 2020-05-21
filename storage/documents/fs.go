package documents

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Shelex/grpc-go-demo/entities"
	"github.com/google/uuid"
)

const (
	localRepo = "storage/documents/local"
)

type FileStorage interface {
	GetEmployeeDocuments(badgeNum int32) ([]string, error)
	GetEmployeeDocument(badgeNum int32, filename string) (entities.Document, error)
	SaveDocument(badgeNum int32, data []byte) error
	RemoveDocument(badgeNum int32, filename string) error
	RemoveEmployeeDocuments(badgeNum int32) error
}

type LocalFS struct {
	by byBadge
}

type byBadge map[int32]byFileName

type byFileName map[string]entities.Document

func NewLocalFS() FileStorage {
	by := make(byBadge)
	return &LocalFS{
		by: by,
	}
}

func (l *LocalFS) GetEmployeeDocuments(badgeNum int32) ([]string, error) {
	if _, ok := l.by[badgeNum]; ok {
		files := make([]string, 0, len(l.by[badgeNum]))
		for filename := range l.by[badgeNum] {
			files = append(files, filename)
		}
		return files, nil
	}
	return nil, fmt.Errorf("no documents found for employee %d", badgeNum)
}

func (l *LocalFS) GetEmployeeDocument(badgeNum int32, filename string) (entities.Document, error) {
	basePath := filepath.Join(localRepo, strconv.FormatInt(int64(badgeNum), 10))
	var empty entities.Document
	if _, ok := l.by[badgeNum]; ok {
		if document, ok := l.by[badgeNum][filename]; ok {
			path := filepath.Join(basePath, filename)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return empty, err
			}
			document.Data = data
			return document, nil
		}
		return empty, fmt.Errorf("document %s for employee %d not found", filename, badgeNum)
	}
	return empty, fmt.Errorf("no documents found for employee %d", badgeNum)
}

func (l *LocalFS) SaveDocument(badgeNum int32, data []byte) error {
	badgePath := strconv.FormatInt(int64(badgeNum), 10)
	basePath := filepath.Join(localRepo, badgePath)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	filename := fmt.Sprintf("%d-%s.png", badgeNum, uuid.New().String())
	path := filepath.Join(basePath, filename)

	if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
		return err
	}

	if _, ok := l.by[badgeNum]; !ok {
		l.by[badgeNum] = make(byFileName)
	}

	document := entities.Document{
		FileName:  filename,
		Data:      []byte{},
		CreatedAt: time.Now().Unix(),
	}

	l.by[badgeNum][filename] = document

	return nil
}

func (l *LocalFS) RemoveDocument(badgeNum int32, filename string) error {

	if _, ok := l.by[badgeNum]; !ok {
		return fmt.Errorf("documents for employee %d not found", badgeNum)
	}

	if _, ok := l.by[badgeNum][filename]; !ok {
		return fmt.Errorf("file %s for employee %d not found", filename, badgeNum)
	}

	path := filepath.Join(localRepo, strconv.FormatInt(int64(badgeNum), 10), filename)

	if err := os.Remove(path); err != nil {
		return err
	}

	delete(l.by[badgeNum], filename)

	return nil
}

func (l *LocalFS) RemoveEmployeeDocuments(badgeNum int32) error {
	if _, ok := l.by[badgeNum]; !ok {
		return fmt.Errorf("documents for employee %d not found", badgeNum)
	}

	path := filepath.Join(localRepo, strconv.FormatInt(int64(badgeNum), 10))

	if err := os.RemoveAll(path); err != nil {
		return err
	}

	delete(l.by, badgeNum)

	return nil
}
