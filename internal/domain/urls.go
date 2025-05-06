package domain

// ShortenRequest is a base struct for shortening url
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenerResponse is a response struct when you request batch urls
type ShortenerBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// URLS is a struct from batch request
type URLS struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
}

// UserURLs is a response struct for user urls
type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ResponseError is a struct for error responses
type ResponseError struct {
	Description string `json:"description"`
}

// URLStorage is a simple map for save urls
type URLStorage map[string]string

// Stats is a struct for stats
type Stats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
