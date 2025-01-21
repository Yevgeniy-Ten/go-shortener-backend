package domain

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenerBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
type URLS struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
}

type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ResponseError struct {
	Description string `json:"description"`
}

type URLStorage map[string]string
