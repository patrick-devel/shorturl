package handlers

import (
	"io"
	"math/big"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sqids/sqids-go"
)

const (
	localhost = "http://localhost:8080"
	minLength = 6
)

var Cache = map[string]string{}

func GenerateHash(url string) (*string, error) {
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

func MakeShortLink(context *gin.Context) {
	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	urlBase, err := url.ParseRequestURI(string(body))
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	urlHashBytes, err := GenerateHash(urlBase.String())
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	shortLink := localhost + "/" + *urlHashBytes
	Cache[*urlHashBytes] = urlBase.String()

	context.String(http.StatusCreated, shortLink)

}

func RedirectShortLink(context *gin.Context) {
	hashURL := context.Param("id")
	baseURL, ok := Cache[hashURL]
	if !ok {
		context.String(http.StatusBadRequest, "link does not exist")
		return
	}

	context.Redirect(http.StatusTemporaryRedirect, baseURL)
}
