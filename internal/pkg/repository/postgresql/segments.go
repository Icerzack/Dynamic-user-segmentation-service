package postgresql

import (
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/repository"
	"context"
	"database/sql"
	"time"
)

type SegmentsRepo struct {
	db db.DBops
}

func NewSegmentsRepo(db db.DBops) *SegmentsRepo {
	return &SegmentsRepo{db: db}
}

func (r *SegmentsRepo) Add(ctx context.Context, segment *repository.Segment) error {
	var title string
	err := r.db.ExecQueryRow(ctx, "INSERT INTO segments(title, description, created_at, updated_at) VALUES ($1, $2, $3, $3) ON CONFLICT (title) DO UPDATE SET title = excluded.title RETURNING title ", segment.Title, segment.Description, time.Now()).Scan(&title)
	return err
}

func (r *SegmentsRepo) GetByTitle(ctx context.Context, title string) (*repository.Segment, error) {
	var segment repository.Segment
	err := r.db.Get(ctx, &segment, "SELECT title, description, created_at, updated_at FROM segments WHERE title = $1", title)
	if err == sql.ErrNoRows {
		return nil, repository.ErrObjectNotFound
	}
	return &segment, err
}

func (r *SegmentsRepo) List(ctx context.Context) ([]*repository.Segment, error) {
	segments := make([]*repository.Segment, 0)
	err := r.db.Select(ctx, &segments, "SELECT title, description, created_at, updated_at FROM segments")
	return segments, err
}

func (r *SegmentsRepo) Update(ctx context.Context, segment *repository.Segment, title string) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE segments SET title = $1, description = $2, updated_at = $3 WHERE title = $4", segment.Title, segment.Description, time.Now(), title)
	return result.RowsAffected() > 0, err
}

func (r *SegmentsRepo) Delete(ctx context.Context, segment *repository.Segment) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM segments WHERE title = $1", segment.Title)
	return result.RowsAffected() > 0, err
}
