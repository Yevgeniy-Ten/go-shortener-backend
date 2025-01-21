package domain

type urlStorage interface {
	Save(value string) (string, error)
	GetURL(shortURL string) string
	SaveBatch(urls []URLS) error
}

type userStorage interface {
	Create() (int, error)
}
type Storage struct {
	URLS urlStorage
	User userStorage
}
