// Code generated by MockGen. DO NOT EDIT.
// Source: author_repo.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "article-service/model"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIAuthorRepository is a mock of IAuthorRepository interface.
type MockIAuthorRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAuthorRepositoryMockRecorder
}

// MockIAuthorRepositoryMockRecorder is the mock recorder for MockIAuthorRepository.
type MockIAuthorRepositoryMockRecorder struct {
	mock *MockIAuthorRepository
}

// NewMockIAuthorRepository creates a new mock instance.
func NewMockIAuthorRepository(ctrl *gomock.Controller) *MockIAuthorRepository {
	mock := &MockIAuthorRepository{ctrl: ctrl}
	mock.recorder = &MockIAuthorRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAuthorRepository) EXPECT() *MockIAuthorRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockIAuthorRepository) Get(ctx context.Context, id uuid.UUID) (*model.Author, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*model.Author)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIAuthorRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIAuthorRepository)(nil).Get), ctx, id)
}
