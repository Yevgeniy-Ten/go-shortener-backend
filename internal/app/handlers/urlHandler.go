package handlers

import (
	"io"
	"net/http"
	"shorter/internal/app/storage"
)

func PostHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Ошибка чтения тела запроса."))
		return
	}
	id := storage.GlobalUrlStorage.Save(string(body))
	respText := "http://localhost:8080/" + id
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(respText))
}
func GetHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[1:]
	url := storage.GlobalUrlStorage.GetUrl(id)
	if url == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Ссылка не найдена."))
		return
	}
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
func URLHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		PostHandler(res, req)
		return
	}
	if req.Method == http.MethodGet {
		GetHandler(res, req)
		return
	}
	res.WriteHeader(http.StatusBadRequest)
}
