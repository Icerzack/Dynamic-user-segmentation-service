package repository

import (
	"context"
	"errors"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type SegmentsRepo interface {
	Add(ctx context.Context, segment *Segment) error
	GetByTitle(ctx context.Context, title string) (*Segment, error)
	List(ctx context.Context) ([]*Segment, error)
	Update(ctx context.Context, segment *Segment, title string) (bool, error)
}

type UsersSegmentsRepo interface {
	Add(ctx context.Context, userSegment *UserSegment) error
	GetByID(ctx context.Context, id int) (*UserSegment, error)
	GetByUserID(ctx context.Context, userID int) ([]*UserSegment, error)
	GetBySegmentTitle(ctx context.Context, segmentTitle string) ([]*UserSegment, error)
	List(ctx context.Context) ([]*UserSegment, error)
	Delete(ctx context.Context, userSegment *UserSegment) (bool, error)
}
