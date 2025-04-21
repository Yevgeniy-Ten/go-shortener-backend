package urls

// DuplicateError is an error when shortURL already exists, custom error
type DuplicateError struct {
	ExistingURL string
	ShortURL    string
}

// Error is a method for DuplicateError
func (de *DuplicateError) Error() string {
	return "ShortURL already exists"
}

// NewDuplicateError is a constructor for DuplicateError
func NewDuplicateError(existingURL, shortURL string) error {
	return &DuplicateError{
		ExistingURL: existingURL,
		ShortURL:    shortURL,
	}
}

// URLIsDeletedError is an error when shortURL is deleted, custom error
type URLIsDeletedError struct {
	ShortURL string
}

// NewURLIsDeletedError is a constructor for URLIsDeletedError
func NewURLIsDeletedError(shortURL string) error {
	return &URLIsDeletedError{
		ShortURL: shortURL,
	}
}

// Error is a method for URLIsDeletedError
func (ud *URLIsDeletedError) Error() string {
	return "ShortURL deleted"
}
