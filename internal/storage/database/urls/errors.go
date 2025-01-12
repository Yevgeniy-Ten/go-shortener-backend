package urls

type DuplicateError struct {
	ExistingURL string
	ShortURL    string
}

func (de *DuplicateError) Error() string {
	return "URL already exists"
}

func NewDuplicateError(existingURL, shortURL string) error {
	return &DuplicateError{
		ExistingURL: existingURL,
		ShortURL:    shortURL,
	}
}
