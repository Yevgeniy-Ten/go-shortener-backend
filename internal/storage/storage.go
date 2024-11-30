package storage

import (
	generateRandomId "shorter/pkg"
)

type ShortURLStorage struct {
	storage map[string]string
}

func NewShortURLStorage() *ShortURLStorage {
	return &ShortURLStorage{
		storage: make(map[string]string),
	}
}

//func NewShortURLStorage() ShortURLStorage {
//	return ShortURLStorage{
//		storage: make(map[string]string),
//	}
//}

func (storage *ShortURLStorage) Save(url string) string {
	newID := generateRandomId.GenerateShortID()
	storage.storage[newID] = url
	return newID
}

//	func (storage *ShortURLStorage) SaveTest(url string) *ShortURLStorage {
//		newID := generateRandomId.GenerateShortID()
//		storage.storage[newID] = url
//		return storage
//	}
func (storage *ShortURLStorage) GetURL(id string) string {
	return storage.storage[id]
}

var GlobalURLStorage = NewShortURLStorage()
