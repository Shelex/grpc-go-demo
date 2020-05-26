package documents

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Shelex/grpc-go-demo/domain/entities"
	"github.com/google/uuid"
)

const (
	localRepo = "storage/documents/local"
)

type FileStorage interface {
	SaveDocument(userID string, filename string, data []byte) (entities.Document, error)
	GetDocument(id string) (entities.Document, error)
	DeleteDocument(id string) error
}

type LocalFS map[string]entities.Document

func NewLocalFS() FileStorage {
	fs := make(LocalFS)
	return &fs
}

func (l *LocalFS) GetDocument(id string) (entities.Document, error) {
	doc, ok := (*l)[id]
	if !ok {
		return doc, fmt.Errorf("document with id %s not found", id)
	}

	ext := filepath.Ext(doc.FileName)

	path := filepath.Join(localRepo, doc.ID+ext)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return doc, fmt.Errorf("document with id %s not found", id)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return doc, err
	}
	doc.Data = data
	return doc, nil
}

func (l *LocalFS) SaveDocument(userID string, filename string, data []byte) (entities.Document, error) {
	var empty entities.Document
	if _, err := os.Stat(localRepo); os.IsNotExist(err) {
		if err := os.MkdirAll(localRepo, os.ModePerm); err != nil {
			return empty, err
		}
	} else if err != nil {
		return empty, err
	}
	extension := filepath.Ext(filename)
	newID, err := uuid.NewUUID()
	if err != nil {
		return empty, err
	}
	documentID := newID.String()

	path := filepath.Join(localRepo, documentID+extension)

	if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
		return empty, err
	}

	if _, ok := (*l)[documentID]; ok {
		return empty, fmt.Errorf("document with id %s already exists", documentID)
	}

	document := entities.Document{
		ID:        documentID,
		FileName:  filename,
		Data:      []byte{},
		CreatedAt: time.Now().Unix(),
	}

	(*l)[documentID] = document

	return document, nil
}

func (l *LocalFS) DeleteDocument(id string) error {
	doc, ok := (*l)[id]
	if !ok {
		return fmt.Errorf("document with id %s not found", id)
	}

	ext := filepath.Ext(doc.FileName)

	path := filepath.Join(localRepo, doc.ID+ext)

	if err := os.Remove(path); err != nil {
		return err
	}

	delete((*l), id)

	return nil
}
