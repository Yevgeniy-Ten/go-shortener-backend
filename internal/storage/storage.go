package storage

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"shorter/internal/domain"
	generateRandomId "shorter/pkg"
	"sync"
)

const (
	FilePerm = os.FileMode(0666)
)

type database interface {
	Save(ctx context.Context, values domain.URLS) error
	GetURL(shortURL string) (string, error)
}

type ShortURLStorage struct {
	storage map[string]string
	mutex   *sync.Mutex
	file    io.WriteCloser
	writer  *bufio.Writer
	scanner *bufio.Scanner
	db      database
}

func NewShortURLStorage(filePath string, db database) (*ShortURLStorage, error) {
	storage := make(map[string]string)
	if filePath == "" {
		return &ShortURLStorage{
			storage: storage,
			mutex:   &sync.Mutex{},
			db:      db,
		}, nil
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, FilePerm)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	parseLine := func(line string) (string, string) {
		return line[:8], line[9:]
	}

	for scanner.Scan() {
		line := scanner.Text()
		id, url := parseLine(line)
		storage[id] = url
	}

	return &ShortURLStorage{
		storage: storage,
		mutex:   &sync.Mutex{},
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: scanner,
		db:      db,
	}, nil
}
func (storage *ShortURLStorage) WriteToFile(newID, url string) error {
	_, err := storage.writer.WriteString(newID + " " + url + "\n")
	if err != nil {
		return err
	}

	err = storage.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
func (storage *ShortURLStorage) Save(url string) (string, error) {
	newID := generateRandomId.GenerateShortID()
	var err error
	if storage.db != nil {
		err = storage.db.Save(context.TODO(), domain.URLS{
			ShortURL: newID,
			URL:      url,
		})
		if err != nil {
			return "", err
		}
	}
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	if storage.file != nil {
		err = storage.WriteToFile(newID, url)
		if err != nil {
			return "", err
		}
	}
	storage.storage[newID] = url
	return newID, nil
}

func (storage *ShortURLStorage) GetURL(id string) string {
	if storage.db != nil {
		fmt.Println(reflect.TypeOf(storage.db), "storage.db", reflect.ValueOf(storage.db))
		url, err := storage.db.GetURL(id)
		if err != nil {
			return ""
		}
		return url
	}
	return storage.storage[id]
}

func (storage *ShortURLStorage) Close() error {
	if storage.file == nil {
		return nil
	}
	return storage.file.Close()
}
