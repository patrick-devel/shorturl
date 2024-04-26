// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/kspopova/GolandProjects/shorturl/internal/handlers/shorturl.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/patrick-devel/shorturl/internal/models"
)

// MockshortService is a mock of shortService interface.
type MockshortService struct {
	ctrl     *gomock.Controller
	recorder *MockshortServiceMockRecorder
}

// MockshortServiceMockRecorder is the mock recorder for MockshortService.
type MockshortServiceMockRecorder struct {
	mock *MockshortService
}

// NewMockshortService creates a new mock instance.
func NewMockshortService(ctrl *gomock.Controller) *MockshortService {
	mock := &MockshortService{ctrl: ctrl}
	mock.recorder = &MockshortServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockshortService) EXPECT() *MockshortServiceMockRecorder {
	return m.recorder
}

// GetOriginalURL mocks base method.
func (m *MockshortService) GetOriginalURL(ctx context.Context, hash string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURL", ctx, hash)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURL indicates an expected call of GetOriginalURL.
func (mr *MockshortServiceMockRecorder) GetOriginalURL(ctx, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURL", reflect.TypeOf((*MockshortService)(nil).GetOriginalURL), ctx, hash)
}

// MakeShortURL mocks base method.
func (m *MockshortService) MakeShortURL(ctx context.Context, originalURL, uid string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeShortURL", ctx, originalURL, uid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeShortURL indicates an expected call of MakeShortURL.
func (mr *MockshortServiceMockRecorder) MakeShortURL(ctx, originalURL, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeShortURL", reflect.TypeOf((*MockshortService)(nil).MakeShortURL), ctx, originalURL, uid)
}

// MakeShortURLs mocks base method.
func (m *MockshortService) MakeShortURLs(ctx context.Context, bulk models.ListRequestBulk) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeShortURLs", ctx, bulk)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeShortURLs indicates an expected call of MakeShortURLs.
func (mr *MockshortServiceMockRecorder) MakeShortURLs(ctx, bulk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeShortURLs", reflect.TypeOf((*MockshortService)(nil).MakeShortURLs), ctx, bulk)
}
