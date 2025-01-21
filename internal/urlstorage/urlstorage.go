package urlstorage

import (
	"errors"
	"shorter/internal/domain"
	generateRandomId "shorter/pkg"
	"sync"
)

type repository interface {
	Save(values domain.URLS, userID int) error
	GetURL(shortURL string) (string, error)
	GetInitialData() (domain.URLStorage, error)
	SaveBatch(values []domain.URLS, userID int) error
	GetUserURLs(userID int) ([]domain.UserURLs, error)
}

type ShortURLStorage struct {
	storage domain.URLStorage
	mutex   *sync.Mutex
	db      repository
}

func (s *ShortURLStorage) GetUserURLs(userID int) ([]domain.UserURLs, error) {
	if s.db != nil {
		return s.db.GetUserURLs(userID)
	}
	return nil, errors.New("not implemented")
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

func (s *ShortURLStorage) Save(url string, userID int) (string, error) {
	newID := generateRandomId.GenerateShortID()
	var err error
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.db != nil {
		err = s.db.Save(domain.URLS{
			CorrelationID: newID,
			URL:           url,
		}, userID)
		if err != nil {
			return "", err
		}
	}
	s.storage[newID] = url
	return newID, nil
}

func (s *ShortURLStorage) GetURL(id string) string {
	if s.db != nil {
		url, err := s.db.GetURL(id)
		if err != nil {
			return ""
		}
		return url
	}
	return s.storage[id]
}
func (s *ShortURLStorage) SaveBatch(urls []domain.URLS, userID int) error {
	_ = s.db.SaveBatch(urls, userID)
	for _, value := range urls {
		s.storage[value.CorrelationID] = value.URL
	}
	return nil
}
