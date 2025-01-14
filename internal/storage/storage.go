package storage

import (
	"errors"
	"shorter/internal/domain"
	generateRandomId "shorter/pkg"
	"sync"
)

type repository interface {
	Save(values domain.URLS) error
	GetURL(shortURL string) (string, error)
	GetInitialData() (domain.Storage, error)
	SaveBatch(values []domain.URLS) error
	Ping() error
}

type ShortURLStorage struct {
	storage domain.Storage
	mutex   *sync.Mutex
	db      repository
}

func New(db repository) *ShortURLStorage {
	storage := make(domain.Storage)
	if db != nil {
		if initialData, err := db.GetInitialData(); err == nil {
			storage = initialData
		}
	}

	return &ShortURLStorage{
		storage: storage,
		mutex:   &sync.Mutex{},
		db:      db,
	}
}

func (storage *ShortURLStorage) Save(url string) (string, error) {
	newID := generateRandomId.GenerateShortID()
	var err error
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	if storage.db != nil {
		err = storage.db.Save(domain.URLS{
			CorrelationID: newID,
			URL:           url,
		})
		if err != nil {
			return "", err
		}
	}
	storage.storage[newID] = url
	return newID, nil
}

func (storage *ShortURLStorage) GetURL(id string) string {
	if storage.db != nil {
		url, err := storage.db.GetURL(id)
		if err != nil {
			return ""
		}
		return url
	}
	return storage.storage[id]
}
func (storage *ShortURLStorage) SaveBatch(urls []domain.URLS) error {
	_ = storage.db.SaveBatch(urls)
	for _, value := range urls {
		storage.storage[value.CorrelationID] = value.URL
	}
	return nil
}
func (storage *ShortURLStorage) Ping() error {
	if storage.db == nil {
		return errors.New("connection is nil")
	}
	return storage.db.Ping()
}
