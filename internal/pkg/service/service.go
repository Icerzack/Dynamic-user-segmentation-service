package service

import (
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/model"
	"context"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service interface {
	InsertSegmentIntoDatabase(ctx context.Context, db db.DBops, request model.SegmentRequest) error
	DeleteSegmentFromDatabase(ctx context.Context, db db.DBops, request model.SegmentRequest) (bool, error)
	ModifyUsersSegmentsInDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest) ([]string, []string, error)
	GetUserSegmentsFromDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest) ([]model.UserSegments, error)
}

type ServiceStub struct {
}

func NewServiceStub() *ServiceStub {
	return &ServiceStub{}
}

func (s *ServiceStub) InsertSegmentIntoDatabase(ctx context.Context, db db.DBops, request model.SegmentRequest) error {
	return nil
}

func (s *ServiceStub) DeleteSegmentFromDatabase(ctx context.Context, db db.DBops, request model.SegmentRequest) (bool, error) {
	return false, nil
}

func (s *ServiceStub) ModifyUsersSegmentsInDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest) ([]string, []string, error) {
	return nil, nil, nil
}
func (s *ServiceStub) GetUserSegmentsFromDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest) ([]model.UserSegments, error) {
	return nil, nil
}
