package db

import (
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/model"
	"avito-backend-internship/internal/pkg/repository"
	"avito-backend-internship/internal/pkg/repository/postgresql"
	"avito-backend-internship/internal/pkg/service/history"
	"context"
	"strconv"
	"time"
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

func (s *PostgresService) ModifyUsersSegmentsInDatabase(ctx context.Context, db db.DBops, request model.UserSegmentRequest, historyService history.Service) ([]string, []string, error) {
	usersSegmentsRepo := postgresql.NewUsersSegmentsRepo(db)
	segmentsRepo := postgresql.NewSegmentsRepo(db)

	addNotExist := make([]string, 0)
	deleteNotExist := make([]string, 0)

	segs, _ := usersSegmentsRepo.GetByUserID(ctx, *request.UserID)
	if request.SegmentsTitlesToAdd != nil {
		for i, segment := range *request.SegmentsTitlesToAdd {
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
				historyService.WriteToFile(*request.UserID, segment, "ADD", time.Now())
				if request.TTL != nil {
					if i < len(*request.TTL) && (*request.TTL)[i] != 0 {
						go s.ttlDetector(ctx, db, *request.UserID, segment, (*request.TTL)[i])
					}
				}
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
			historyService.WriteToFile(*request.UserID, segment, "DEL", time.Now())
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

func (s *PostgresService) ttlDetector(ctx context.Context, db db.DBops, userID int, segmentTitle string, ttlSeconds int) {
	usersSegmentsRepo := postgresql.NewUsersSegmentsRepo(db)

	stringTime, _ := time.ParseDuration(strconv.Itoa(ttlSeconds) + "s")
	now := time.NewTimer(stringTime)
	<-now.C
	_, err := usersSegmentsRepo.Delete(ctx, &repository.UserSegment{
		UserID:       userID,
		SegmentTitle: segmentTitle,
	})
	if err != nil {
		return
	}
}
