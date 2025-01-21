package urlstorage

import (
	"shorter/internal/domain"
	generateRandomId "shorter/pkg"
	"sync"
)

type repository interface {
	Save(values domain.URLS, userID int) error
	GetURL(shortURL string) (string, error)
	GetInitialData() (domain.URLStorage, error)
	SaveBatch(values []domain.URLS, userID int) error
}

type ShortURLStorage struct {
	storage domain.URLStorage
	mutex   *sync.Mutex
	db      repository
}

func New(db repository) *ShortURLStorage {
	storage := make(domain.URLStorage)
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
		}, 0)
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
func (storage *ShortURLStorage) SaveBatch(urls []domain.URLS, userID int) error {
	_ = storage.db.SaveBatch(urls, userID)
	for _, value := range urls {
		storage.storage[value.CorrelationID] = value.URL
	}
	return nil
}

func (storage *ShortURLStorage) SaveWithUser(url string, userID int) (string, error) {
	newID := generateRandomId.GenerateShortID()
	var err error
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	if storage.db != nil {
		err = storage.db.Save(domain.URLS{
			CorrelationID: newID,
			URL:           url,
		}, userID)
		if err != nil {
			return "", err
		}
	}
	storage.storage[newID] = url
	return newID, nil
}
