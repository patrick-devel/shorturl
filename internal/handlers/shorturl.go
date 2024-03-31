package handlers

import (
	"io"
	"math/big"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sqids/sqids-go"

	"github.com/patrick-devel/shorturl/config"
	"github.com/patrick-devel/shorturl/internal/models"
)

const minLength = 6

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

func MakeShortLinkHandler(c *config.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
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

		shortLink := c.BaseURL.String() + "/" + *urlHashBytes
		Cache[*urlHashBytes] = urlBase.String()

		context.String(http.StatusCreated, shortLink)
	}
}

func RedirectShortLinkHandler(context *gin.Context) {
	hashURL := context.Param("id")
	baseURL, ok := Cache[hashURL]
	if !ok {
		context.String(http.StatusBadRequest, "link does not exist")
		return
	}

	context.Redirect(http.StatusTemporaryRedirect, baseURL)
}

func MakeShortURLJSONHandler(c *config.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		var request models.Request

		if err := context.BindJSON(&request); err != nil {
			context.String(http.StatusBadRequest, err.Error())
		}

		urlHashBytes, err := GenerateHash(request.URL.String())
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}

		Cache[*urlHashBytes] = request.URL.String()
		resp := models.Response{Result: c.BaseURL.String() + "/" + *urlHashBytes}
		context.JSON(http.StatusCreated, resp)
	}
}
