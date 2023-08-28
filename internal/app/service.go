package app

import (
	"avito-backend-internship/internal/pkg/model"
	"avito-backend-internship/internal/pkg/repository"
	"avito-backend-internship/internal/pkg/repository/postgresql"
)

func (s *Server) insertSegmentIntoDatabase(request model.SegmentRequest) error {
	segmentsRepo := postgresql.NewSegmentsRepo(s.db)
	if request.Description == nil {
		err := segmentsRepo.Add(s.ctx, &repository.Segment{
			Title:       request.Title,
			Description: "No description provided.",
		})
		return err
	}
	err := segmentsRepo.Add(s.ctx, &repository.Segment{
		Title:       request.Title,
		Description: *request.Description,
	})

	return err
}

func (s *Server) deleteSegmentFromDatabase(request model.SegmentRequest) (bool, error) {
	segmentsRepo := postgresql.NewSegmentsRepo(s.db)
	ok, err := segmentsRepo.Delete(s.ctx, &repository.Segment{
		Title: request.Title,
	})
	if ok {
		return true, err
	}
	return false, err
}

func (s *Server) modifyUsersSegmentsInDatabase(request model.UserSegmentRequest) ([]string, []string, error) {
	usersSegmentsRepo := postgresql.NewUsersSegmentsRepo(s.db)
	segmentsRepo := postgresql.NewSegmentsRepo(s.db)

	addNotExist := make([]string, 0)
	deleteNotExist := make([]string, 0)

	segs, _ := usersSegmentsRepo.GetByUserID(s.ctx, request.UserID)
	for _, segment := range request.SegmentsTitlesToAdd {
		_, err := segmentsRepo.GetByTitle(s.ctx, segment)
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
			err = usersSegmentsRepo.Add(s.ctx, &repository.UserSegment{
				UserID:       request.UserID,
				SegmentTitle: segment,
			})
		}
		if err != nil {
			return nil, nil, err
		}
	}

	for _, segment := range request.SegmentsTitlesToDelete {
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
		_, err := usersSegmentsRepo.Delete(s.ctx, &repository.UserSegment{
			UserID:       request.UserID,
			SegmentTitle: segment,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	return addNotExist, deleteNotExist, nil
}

func (s *Server) getUserSegmentsFromDatabase(request model.UserSegmentRequest) ([]model.UserSegments, error) {
	usersSegmentsRepo := postgresql.NewUsersSegmentsRepo(s.db)
	segmentsRepo := postgresql.NewSegmentsRepo(s.db)

	segments, err := usersSegmentsRepo.GetByUserID(s.ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	array := make([]model.UserSegments, len(segments))

	for id, segment := range segments {
		seg, err := segmentsRepo.GetByTitle(s.ctx, segment.SegmentTitle)
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
