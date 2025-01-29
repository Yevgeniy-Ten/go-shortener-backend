package urlstorage

import (
	"context"
	"errors"
	"shorter/internal/domain"
	generateRandomId "shorter/pkg"
	"sync"
)

type repository interface {
	Save(values domain.URLS, userID int) error
	GetURL(shortURL string) (string, error)
	GetInitialData() (domain.URLStorage, error)
	SaveBatch(ctx context.Context, values []domain.URLS, userID int) error
	GetUserURLs(userID int, serverAdr string) ([]domain.UserURLs, error)
	DeleteURLs(correlationIDS []string, userID int) error
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
func (s *ShortURLStorage) DeleteURLs(correlationIDS []string, userID int) error {
	if s.db != nil {
		return s.db.DeleteURLs(correlationIDS, userID)
	}
	return errors.New("not implemented")
}

func (s *ShortURLStorage) GetUserURLs(userID int, serverAdr string) ([]domain.UserURLs, error) {
	if s.db != nil {
		return s.db.GetUserURLs(userID, serverAdr)
	}
	return nil, errors.New("not implemented")
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

func (s *ShortURLStorage) GetURL(id string) (string, error) {
	if s.db != nil {
		url, err := s.db.GetURL(id)
		if err != nil {
			return "", err
		}
		return url, nil
	}
	return s.storage[id], nil
}
func (s *ShortURLStorage) SaveBatch(urls []domain.URLS, userID int) error {
	ctx := context.TODO()
	_ = s.db.SaveBatch(ctx, urls, userID)
	for _, value := range urls {
		s.storage[value.CorrelationID] = value.URL
	}
	return nil
}
