package domain

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenerBatchResponse struct {
	URLId string `json:"correlation_id"`
	URL   string `json:"short_url"`
}
type URLS struct {
	URLId string `json:"correlation_id"`
	URL   string `json:"original_url"`
}

type Storage map[string]string
