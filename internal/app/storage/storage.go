package storage

import generateRandomId "shorter/internal/app/lib"

type ShortUrlStorage struct {
	storage map[string]string
}

func NewShortUrlStorage() *ShortUrlStorage {
	return &ShortUrlStorage{
		storage: make(map[string]string),
	}
}

func (storage *ShortUrlStorage) Save(url string) string {
	newId := generateRandomId.GenerateShortId()
	storage.storage[newId] = url
	return newId
}

func (storage *ShortUrlStorage) GetUrl(id string) string {
	return storage.storage[id]
}

var GlobalUrlStorage = NewShortUrlStorage()
