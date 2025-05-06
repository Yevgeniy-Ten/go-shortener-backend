package domain

type urlStorage interface {
	Save(value string, userID int) (string, error)
	GetURL(shortURL string) (string, error)
	SaveBatch(urls []URLS, userID int) error
	GetUserURLs(userID int, serverAdr string) ([]UserURLs, error)
	DeleteURLs(correlationIDS []string, userID int) error
	GetStats() (*Stats, error)
}

type userStorage interface {
	Create() (int, error)
}

// Storage struct in project
//
//go:generate mockgen -source=internal/domain/storage.go -destination=internal/handlers/mocks/mock_storage.go -package=handlers_mocks Storage
type Storage struct {
	URLS urlStorage
	User userStorage
}
