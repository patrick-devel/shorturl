package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patrick-devel/shorturl/config"
	"github.com/patrick-devel/shorturl/internal/handlers"
)

var defaultConfig = config.NewConfigBuilder().Build()

func TestMakeShortLinkHandler(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/", handlers.MakeShortLinkHandler(&defaultConfig))
	router.HandleMethodNotAllowed = true

	tests := []struct {
		name    string
		method  string
		body    io.Reader
		expCode int
	}{
		{
			name:    "OK",
			method:  http.MethodPost,
			body:    strings.NewReader("https://practicum.yandex.ru/"),
			expCode: http.StatusCreated,
		},
		{
			name:    "MethodNotAllowed",
			method:  http.MethodPatch,
			body:    strings.NewReader("https://practicum.yandex.ru/"),
			expCode: http.StatusMethodNotAllowed,
		},
		{
			name:    "BadRequest",
			method:  http.MethodPost,
			body:    strings.NewReader("https://\\]]practicum.yandex.ru/"),
			expCode: http.StatusBadRequest,
		},
		{
			name:    "BadRequestBodyEmpty",
			method:  http.MethodPost,
			body:    strings.NewReader(""),
			expCode: http.StatusBadRequest,
		},
		{
			name:    "BadRequestBodyIncorrect",
			method:  http.MethodPost,
			body:    strings.NewReader("make short url pls"),
			expCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(testcase.method, "/", testcase.body)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, testcase.expCode, resp.StatusCode)
			assert.NotEmpty(t, body)
		})
	}
}

func TestRedirectShortLinkHandler(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:id", handlers.RedirectShortLinkHandler)
	router.HandleMethodNotAllowed = true

	baseURL := "https://practicum.yandex.ru/"
	hash, err := handlers.GenerateHash(baseURL)
	require.NoError(t, err)

	handlers.Cache[*hash] = baseURL

	tests := []struct {
		name    string
		method  string
		hash    string
		expCode int
	}{
		{
			name:    "OK",
			method:  http.MethodGet,
			hash:    *hash,
			expCode: http.StatusTemporaryRedirect,
		},
		{
			name:    "MethodNotAllowed",
			method:  http.MethodPatch,
			hash:    *hash,
			expCode: http.StatusMethodNotAllowed,
		},
		{
			name:    "BadRequest",
			method:  http.MethodGet,
			hash:    "not_exist",
			expCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(testcase.method, fmt.Sprintf("/%s", testcase.hash), http.NoBody)
			req.SetPathValue("id", testcase.hash)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			assert.Equal(t, testcase.expCode, resp.StatusCode)
			if resp.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, baseURL, resp.Header.Get("Location"))
			}
		})
	}
}

func TestMakeShortLinkJSONHandler(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/", handlers.MakeShortURLJSONHandler(&defaultConfig))
	router.HandleMethodNotAllowed = true

	tests := []struct {
		name    string
		method  string
		body    io.Reader
		expCode int
	}{
		{
			name:    "OK",
			method:  http.MethodPost,
			body:    strings.NewReader(`{"url": "https://practicum.yandex.ru/"}`),
			expCode: http.StatusCreated,
		},
		{
			name:    "MethodNotAllowed",
			method:  http.MethodPatch,
			body:    strings.NewReader(`{"url": "https://practicum.yandex.ru/"}`),
			expCode: http.StatusMethodNotAllowed,
		},
		{
			name:    "BadRequest",
			method:  http.MethodPost,
			body:    strings.NewReader(`{"url": "https://\\]]practicum.yandex.ru/"}`),
			expCode: http.StatusBadRequest,
		},
		{
			name:    "BadRequestBodyEmpty",
			method:  http.MethodPost,
			body:    strings.NewReader(`{}`),
			expCode: http.StatusBadRequest,
		},
		{
			name:    "BadRequestBodyIncorrect",
			method:  http.MethodPost,
			body:    strings.NewReader(`make short url pls`),
			expCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(testcase.method, "/", testcase.body)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, testcase.expCode, resp.StatusCode)

			assert.NotEmpty(t, body)
		})
	}
}
