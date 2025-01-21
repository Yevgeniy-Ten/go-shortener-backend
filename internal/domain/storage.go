package domain

type urlStorage interface {
	Save(value string, userID int) (string, error)
	GetURL(shortURL string) string
	SaveBatch(urls []URLS, userID int) error
}

type userStorage interface {
	Create() (int, error)
}
type Storage struct {
	URLS urlStorage
	User userStorage
}
