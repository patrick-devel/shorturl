package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/patrick-devel/shorturl/internal/handlers"
	mockhandlers "github.com/patrick-devel/shorturl/internal/handlers/mocks"
)

func TestDeleteShortUrls_400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockhandlers.NewMockserviceDeleter(ctrl)

	router := gin.Default()
	router.POST("/delete", handlers.DeleteShortUrls(mockService))

	req, _ := http.NewRequest(http.MethodPost, "/delete", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDeleteShortUrls_202(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockhandlers.NewMockserviceDeleter(ctrl)

	router := gin.Default()
	router.POST("/delete", handlers.DeleteShortUrls(mockService))

	shortUrls := []string{"short1", "short2"}
	mockService.EXPECT().DeleteShortURL(gomock.Any(), shortUrls).Return(nil).AnyTimes()

	reqBody := `["short1", "short2"]`
	req, _ := http.NewRequest(http.MethodPost, "/delete", strings.NewReader(reqBody))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)
}

func TestDeleteShortUrls_202_error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockhandlers.NewMockserviceDeleter(ctrl)

	router := gin.Default()
	router.POST("/delete", handlers.DeleteShortUrls(mockService))

	shortUrls := []string{"short1", "short2"}
	mockService.EXPECT().DeleteShortURL(gomock.Any(), shortUrls).Return(fmt.Errorf("error")).AnyTimes()

	reqBody := `["short1", "short2"]`
	req, _ := http.NewRequest(http.MethodPost, "/delete", strings.NewReader(reqBody))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)
}
