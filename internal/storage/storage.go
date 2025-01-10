package storage

import (
	"shorter/internal/domain"
	generateRandomId "shorter/pkg"
	"sync"
)

type database interface {
	Save(values domain.URLS) error
	GetURL(shortURL string) (string, error)
	GetInitialData() (domain.Storage, error)
	SaveBatch(values []domain.URLS) error
}

type ShortURLStorage struct {
	storage domain.Storage
	mutex   *sync.Mutex
	db      database
}

func NewShortURLStorage(db database) *ShortURLStorage {
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
	if storage.db != nil {
		err = storage.db.Save(domain.URLS{
			URLId: newID,
			URL:   url,
		})
		if err != nil {
			return "", err
		}
	}
	defer storage.mutex.Unlock()
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
		storage.storage[value.URLId] = value.URL
	}
	return nil
}
