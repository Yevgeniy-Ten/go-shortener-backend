package storage

import generateRandomId "shorter/internal/app/lib"

type ShortURLStorage struct {
	storage map[string]string
}

func NewShortURLStorage() *ShortURLStorage {
	return &ShortURLStorage{
		storage: make(map[string]string),
	}
}

func (storage *ShortURLStorage) Save(url string) string {
	newID := generateRandomId.GenerateShortID()
	storage.storage[newID] = url
	return newID
}

func (storage *ShortURLStorage) GetURL(id string) string {
	return storage.storage[id]
}

var GlobalURLStorage = NewShortURLStorage()
