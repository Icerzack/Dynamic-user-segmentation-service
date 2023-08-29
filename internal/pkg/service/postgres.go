package service

import (
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/model"
	"avito-backend-internship/internal/pkg/repository"
	"avito-backend-internship/internal/pkg/repository/postgresql"
	"context"
)

type PostgresService struct{}

func NewPostgresService() *PostgresService {
	return &PostgresService{}
}

func (s *PostgresService) InsertSegmentIntoDatabase(ctx context.Context, db db.DBops, request model.SegmentRequest) error {
	segmentsRepo := postgresql.NewSegmentsRepo(db)
	if request.Description == nil {
		err := segmentsRepo.Add(ctx, &repository.Segment{
			Title:       *request.Title,
			Description: "No description provided.",
		})
		return err
	}
	err := segmentsRepo.Add(ctx, &repository.Segment{
		Title:       *request.Title,
		Description: *request.Description,
	})

	return err
}

func (s *PostgresService) DeleteSegmentFromDatabase(ctx context.Context, db db.DBops, request model.SegmentRequest) (bool, error) {
	segmentsRepo := postgresql.NewSegmentsRepo(db)
	ok, err := segmentsRepo.Delete(ctx, &repository.Segment{
		Title: *request.Title,
	})
	if ok {
		return true, err
	}
	return false, err
}

func (s *PostgresService) ModifyUsersSegmentsInDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest) ([]string, []string, error) {
	usersSegmentsRepo := postgresql.NewUsersSegmentsRepo(db)
	segmentsRepo := postgresql.NewSegmentsRepo(db)

	addNotExist := make([]string, 0)
	deleteNotExist := make([]string, 0)

	segs, _ := usersSegmentsRepo.GetByUserID(ctx, *request.UserID)
	if request.SegmentsTitlesToAdd != nil {
		for _, segment := range *request.SegmentsTitlesToAdd {
			_, err := segmentsRepo.GetByTitle(ctx, segment)
			if err != nil {
				addNotExist = append(addNotExist, segment)
				continue
			}
			var alreadyExists bool
			for _, s := range segs {
				if s.SegmentTitle == segment {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				err = usersSegmentsRepo.Add(ctx, &repository.UserSegment{
					UserID:       *request.UserID,
					SegmentTitle: segment,
				})
			}
			if err != nil {
				return nil, nil, err
			}
		}
	}

	if request.SegmentsTitlesToDelete != nil {
		for _, segment := range *request.SegmentsTitlesToDelete {
			var exists bool
			for _, s := range segs {
				if s.SegmentTitle == segment {
					exists = true
					break
				}
			}
			if !exists {
				deleteNotExist = append(deleteNotExist, segment)
				continue
			}
			_, err := usersSegmentsRepo.Delete(ctx, &repository.UserSegment{
				UserID:       *request.UserID,
				SegmentTitle: segment,
			})
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return addNotExist, deleteNotExist, nil
}

func (s *PostgresService) GetUserSegmentsFromDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest) ([]model.UserSegments, error) {
	usersSegmentsRepo := postgresql.NewUsersSegmentsRepo(db)
	segmentsRepo := postgresql.NewSegmentsRepo(db)

	segments, err := usersSegmentsRepo.GetByUserID(ctx, *request.UserID)
	if err != nil {
		return nil, err
	}

	array := make([]model.UserSegments, len(segments))

	for id, segment := range segments {
		seg, err := segmentsRepo.GetByTitle(ctx, segment.SegmentTitle)
		if err != nil {
			return nil, err
		}
		array[id] = model.UserSegments{
			Title:       seg.Title,
			Description: seg.Description,
		}
	}

	return array, nil
}
