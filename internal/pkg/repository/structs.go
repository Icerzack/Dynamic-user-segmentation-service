package repository

import "time"

type Segment struct {
	Title     string    `db:"title"`
	Desc      string    `db:"desc"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserSegment struct {
	ID           int64     `db:"id"`
	UserID       int       `db:"user_id"`
	SegmentTitle string    `db:"seg_title"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
