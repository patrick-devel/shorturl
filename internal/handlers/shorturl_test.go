package handlers_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patrick-devel/shorturl/internal/handlers"
	"github.com/patrick-devel/shorturl/internal/mocks"
	"github.com/patrick-devel/shorturl/internal/models"
)

func TestMakeShortLinkHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockshortService(ctrl)

	router := gin.Default()
	router.POST("/", handlers.MakeShortLinkHandler(mockService))
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
			mockService.EXPECT().MakeShortURL(gomock.Any(), gomock.Any(), gomock.Any()).Return("http://localhost/test", nil)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockshortService(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:id", handlers.RedirectShortLinkHandler(mockService))
	router.HandleMethodNotAllowed = true

	baseURL := "https://practicum.yandex.ru/"

	tests := []struct {
		name     string
		method   string
		hash     string
		expCode  int
		mockExec func()
	}{
		{
			name:    "OK",
			method:  http.MethodGet,
			hash:    "dkadwda",
			expCode: http.StatusTemporaryRedirect,
			mockExec: func() {
				mockService.EXPECT().GetOriginalURL(gomock.Any(), gomock.Any()).Return(baseURL, nil)
			},
		},
		{
			name:    "MethodNotAllowed",
			method:  http.MethodPatch,
			hash:    "dkadwda",
			expCode: http.StatusMethodNotAllowed,
			mockExec: func() {
				mockService.EXPECT().GetOriginalURL(gomock.Any(), gomock.Any()).Return(baseURL, nil)
			},
		},
		{
			name:    "NotFound",
			method:  http.MethodGet,
			hash:    "not_exist",
			expCode: http.StatusNotFound,
			mockExec: func() {
				mockService.EXPECT().GetOriginalURL(gomock.Any(), gomock.Any()).Return("", errors.New("not found link"))
			},
		},
	}

	for _, tc := range tests {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(testcase.method, fmt.Sprintf("/%s", testcase.hash), http.NoBody)
			req.SetPathValue("id", testcase.hash)

			testcase.mockExec()

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockshortService(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/", handlers.MakeShortURLJSONHandler(mockService))
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

			mockService.EXPECT().
				MakeShortURL(gomock.Any(), gomock.Any(), gomock.Any()).
				Return("http://localhost/123sda", nil)

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

func TestMakeShortLinkBulk(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockshortService(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/", handlers.MakeShortURLBulk(mockService))
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
			body:    strings.NewReader(`[{"correlation_id": "2fb839d2-5bcf-45d1-8c76-64086f726153", "original_url": "https://practicum.yandex.ru/"}]`),
			expCode: http.StatusCreated,
		},
		{
			name:    "MethodNotAllowed",
			method:  http.MethodPatch,
			body:    strings.NewReader(`[{"correlation_id": "2fb839d2-5bcf-45d1-8c76-64086f726153", "original_url": "https://practicum.yandex.ru/"}]`),
			expCode: http.StatusMethodNotAllowed,
		},
		{
			name:    "BadRequest",
			method:  http.MethodPost,
			body:    strings.NewReader(`[{"correlation_id": "2fb839d2-5bcf-45d1-8c76-64086f726153", "original_url": "httpas://sdad\\\practicum.yandex.ru/"}]`),
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
			body:    strings.NewReader(`[]`),
			expCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(testcase.method, "/", testcase.body)
			recorder := httptest.NewRecorder()

			mockService.EXPECT().
				MakeShortURLs(gomock.Any(), gomock.Any()).
				Return([]models.Event{{
					Hash: "dadad", UUID: "2fb839d2-5bcf-45d1-8c76-64086f726153",
					OriginalURL: "https://practicum.yandex.ru/",
					ShortURL:    "http://localhost/123sda"}},
					nil,
				)

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
