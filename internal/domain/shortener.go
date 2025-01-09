package domain

type ShortenRequest struct {
	URL string `json:"url"`
}

type URLS struct {
	ShortURL string
	URL      string
}
