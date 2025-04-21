package urlstorage_test

import (
	"shorter/internal/domain"
	"shorter/internal/urlstorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSave_NoDB(t *testing.T) {
	storage := urlstorage.New(nil)

	url := "https://example.com"
	userID := 1

	shortURL, err := storage.Save(url, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, shortURL)
}

func TestGetURL_NoDB(t *testing.T) {
	storage := urlstorage.New(nil)

	expectedURL := "https://example.com"
	shortURL, err := storage.Save(expectedURL, 1)
	assert.NoError(t, err)
	result, err := storage.GetURL(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, result)
}

func TestSaveBatch(t *testing.T) {
	storage := urlstorage.New(nil)

	userID := 1
	urls := []domain.URLS{
		{CorrelationID: "id1", URL: "https://site1.com"},
		{CorrelationID: "id2", URL: "https://site2.com"},
	}

	err := storage.SaveBatch(urls, userID)
	assert.NoError(t, err)
}

func TestDeleteURLs(t *testing.T) {
	storage := urlstorage.New(nil)

	userID := 1
	ids := []string{"id1", "id2"}

	err := storage.DeleteURLs(ids, userID)
	assert.Error(t, err)
	assert.Equal(t, "not implemented", err.Error())
}

func TestGetUserURLs(t *testing.T) {
	storage := urlstorage.New(nil)

	userID := 1
	serverAdr := "http://localhost"

	_, err := storage.GetUserURLs(userID, serverAdr)
	assert.Error(t, err)
	assert.Equal(t, "not implemented", err.Error())
}
