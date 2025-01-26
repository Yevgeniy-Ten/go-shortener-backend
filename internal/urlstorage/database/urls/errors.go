package urls

type DuplicateError struct {
	ExistingURL string
	ShortURL    string
}

func (de *DuplicateError) Error() string {
	return "ShortURL already exists"
}

func NewDuplicateError(existingURL, shortURL string) error {
	return &DuplicateError{
		ExistingURL: existingURL,
		ShortURL:    shortURL,
	}
}

type URLIsDeletedError struct {
	ShortURL string
}

func NewURLIsDeletedError(shortURL string) error {
	return &URLIsDeletedError{
		ShortURL: shortURL,
	}
}

func (ud *URLIsDeletedError) Error() string {
	return "ShortURL deleted"
}
