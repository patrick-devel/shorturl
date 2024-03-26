package handlers

import (
	"io"
	"math/big"
	"net/http"
	"net/url"

	"github.com/sqids/sqids-go"
)

const (
	localhost = "http://localhost:8080"
	minLength = 6
)

var cache = map[string]string{}

func generateHash(url string) (*string, error) {
	generatedNumber := new(big.Int).SetBytes([]byte(url)).Uint64()
	s, err := sqids.New(sqids.Options{MinLength: minLength})
	if err != nil {
		return nil, err
	}

	id, err := s.Encode([]uint64{generatedNumber})
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func MakeShortLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	urlBase, err := url.ParseRequestURI(string(body))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	urlHashBytes, err := generateHash(urlBase.String())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	shortLink := localhost + "/" + *urlHashBytes
	cache[*urlHashBytes] = urlBase.String()

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortLink))
}

func RedirectShortLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hashURL := req.PathValue("id")
	baseURL, ok := cache[hashURL]
	if !ok {
		http.Error(res, "link does not exist", http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, baseURL, http.StatusTemporaryRedirect)
}
