package postgresql

import (
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/repository"
	"context"
	"database/sql"
	"time"
)

type UsersSegmentsRepo struct {
	db db.DBops
}

func NewUsersSegmentsRepo(db db.DBops) *UsersSegmentsRepo {
	return &UsersSegmentsRepo{db: db}
}

func (r *UsersSegmentsRepo) Add(ctx context.Context, userSegment *repository.UserSegment) error {
	var id int
	err := r.db.ExecQueryRow(ctx, `INSERT INTO users_segments(user_id, seg_title, created_at, updated_at) VALUES ($1, $2, $3, $3) RETURNING id`, userSegment.UserID, userSegment.SegmentTitle, time.Now()).Scan(&id)
	return err
}

func (r *UsersSegmentsRepo) GetByID(ctx context.Context, id int) (*repository.UserSegment, error) {
	var userSegment repository.UserSegment
	err := r.db.Get(ctx, &userSegment, "SELECT id, user_id, seg_title, created_at, updated_at FROM users_segments WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, repository.ErrObjectNotFound
	}
	return &userSegment, err
}

func (r *UsersSegmentsRepo) GetByUserID(ctx context.Context, userID int) ([]*repository.UserSegment, error) {
	userSegments := make([]*repository.UserSegment, 0)
	err := r.db.Select(ctx, &userSegments, "SELECT id, user_id, seg_title, created_at, updated_at FROM users_segments WHERE user_id = $1", userID)
	return userSegments, err
}

func (r *UsersSegmentsRepo) GetBySegmentTitle(ctx context.Context, segmentTitle string) ([]*repository.UserSegment, error) {
	usersSegment := make([]*repository.UserSegment, 0)
	err := r.db.Select(ctx, &usersSegment, "SELECT id, user_id, seg_title, created_at, updated_at FROM users_segments WHERE seg_title = $1", segmentTitle)
	return usersSegment, err
}

func (r *UsersSegmentsRepo) List(ctx context.Context) ([]*repository.UserSegment, error) {
	usersSegments := make([]*repository.UserSegment, 0)
	err := r.db.Select(ctx, &usersSegments, "SELECT id, user_id, seg_title, created_at, updated_at FROM users_segments")
	return usersSegments, err
}

func (r *UsersSegmentsRepo) Delete(ctx context.Context, userSegment *repository.UserSegment) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM users_segments WHERE user_id = $1 AND seg_title = $2", userSegment.UserID, userSegment.SegmentTitle)
	return result.RowsAffected() > 0, err
}
