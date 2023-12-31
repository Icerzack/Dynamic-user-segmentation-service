// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	db "avito-backend-internship/internal/pkg/db"
	model "avito-backend-internship/internal/pkg/model"
	history "avito-backend-internship/internal/pkg/service/history"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// DeleteSegmentFromDatabase mocks base method.
func (m *MockService) DeleteSegmentFromDatabase(ctx context.Context, db db.DBops, request model.SegmentRequest) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSegmentFromDatabase", ctx, db, request)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSegmentFromDatabase indicates an expected call of DeleteSegmentFromDatabase.
func (mr *MockServiceMockRecorder) DeleteSegmentFromDatabase(ctx, db, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSegmentFromDatabase", reflect.TypeOf((*MockService)(nil).DeleteSegmentFromDatabase), ctx, db, request)
}

// GetUserSegmentsFromDatabase mocks base method.
func (m *MockService) GetUserSegmentsFromDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest) ([]model.UserSegments, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSegmentsFromDatabase", ctx, db, request)
	ret0, _ := ret[0].([]model.UserSegments)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSegmentsFromDatabase indicates an expected call of GetUserSegmentsFromDatabase.
func (mr *MockServiceMockRecorder) GetUserSegmentsFromDatabase(ctx, db, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSegmentsFromDatabase", reflect.TypeOf((*MockService)(nil).GetUserSegmentsFromDatabase), ctx, db, request)
}

// InsertSegmentIntoDatabase mocks base method.
func (m *MockService) InsertSegmentIntoDatabase(ctx context.Context, db db.DBops, request model.SegmentRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertSegmentIntoDatabase", ctx, db, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertSegmentIntoDatabase indicates an expected call of InsertSegmentIntoDatabase.
func (mr *MockServiceMockRecorder) InsertSegmentIntoDatabase(ctx, db, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertSegmentIntoDatabase", reflect.TypeOf((*MockService)(nil).InsertSegmentIntoDatabase), ctx, db, request)
}

// ModifyUsersSegmentsInDatabase mocks base method.
func (m *MockService) ModifyUsersSegmentsInDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest, historyService history.Service) ([]string, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModifyUsersSegmentsInDatabase", ctx, db, request, historyService)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ModifyUsersSegmentsInDatabase indicates an expected call of ModifyUsersSegmentsInDatabase.
func (mr *MockServiceMockRecorder) ModifyUsersSegmentsInDatabase(ctx, db, request, historyService interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModifyUsersSegmentsInDatabase", reflect.TypeOf((*MockService)(nil).ModifyUsersSegmentsInDatabase), ctx, db, request, historyService)
}
