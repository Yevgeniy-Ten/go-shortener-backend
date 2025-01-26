package domain

type urlStorage interface {
	Save(value string, userID int) (string, error)
	GetURL(shortURL string) (string, error)
	SaveBatch(urls []URLS, userID int) error
	GetUserURLs(userID int, serverAdr string) ([]UserURLs, error)
	DeleteURLs(correlationIDS []string, userID int) error
}

type userStorage interface {
	Create() (int, error)
}
type Storage struct {
	URLS urlStorage
	User userStorage
}
