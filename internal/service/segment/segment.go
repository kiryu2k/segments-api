package segment

import (
	"context"
	"sync"
	"time"

	"github.com/kiryu-dev/segments-api/internal/model"
	"github.com/kiryu-dev/segments-api/pkg/util/selector"
)

type segmentRepository interface {
	Create(context.Context, string) error
	Delete(context.Context, string) error
	DeleteByTTL(context.Context) ([]*model.UserSegment, error)
	GetUsersBySegment(context.Context, string) ([]uint64, error)
}

type userRepository interface {
	GetAll(context.Context) ([]uint64, error)
	AddSegment(context.Context, *model.UserSegment) error
}

type logsRepository interface {
	Write(context.Context, *model.UserLog) error
}

type Service struct {
	segment segmentRepository
	user    userRepository
	logs    logsRepository
}

type userError struct {
	id  uint64
	err error
}

func New(segment segmentRepository, user userRepository, logs logsRepository) *Service {
	return &Service{segment, user, logs}
}

func (s *Service) Create(ctx context.Context, slug string, percentage float64) ([]uint64, error) {
	err := s.segment.Create(ctx, slug)
	if percentage == 0 || err != nil {
		return nil, err
	}
	users, err := s.user.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	count := len(users)
	if percentage != 100 {
		count = int(percentage / 100. * float64(count))
		if count == 0 {
			return nil, nil
		}
		users, err = selector.Select(users, count)
		if err != nil {
			return nil, err
		}
	}
	result := make([]uint64, 0)
	for e := range s.addSegmentToUsers(ctx, users, slug) {
		if e.err == nil {
			result = append(result, e.id)
		}
	}
	return result, nil
}

func (s *Service) addSegmentToUsers(ctx context.Context, users []uint64, slug string) <-chan *userError {
	var (
		wg  = &sync.WaitGroup{}
		out = make(chan *userError)
	)
	wg.Add(len(users))
	for _, user := range users {
		go func(ctx context.Context, userID uint64) {
			defer wg.Done()
			err := s.user.AddSegment(ctx, &model.UserSegment{
				UserID: userID,
				Slug:   slug,
			})
			if err == nil {
				_ = s.logs.Write(ctx, &model.UserLog{
					UserID:      userID,
					Slug:        slug,
					Operation:   model.AddOp.String(),
					RequestTime: time.Now(),
				})
			}
			out <- &userError{
				id:  userID,
				err: err,
			}
		}(ctx, user)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (s *Service) Delete(ctx context.Context, slug string) error {
	users, _ := s.segment.GetUsersBySegment(ctx, slug)
	if err := s.segment.Delete(ctx, slug); err != nil {
		return err
	}
	for _, id := range users {
		_ = s.logs.Write(ctx, &model.UserLog{
			UserID:      id,
			Slug:        slug,
			Operation:   model.DeleteOp.String(),
			RequestTime: time.Now(),
		})
	}
	return nil
}

func (s *Service) DeleteByTTL() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	segments, err := s.segment.DeleteByTTL(ctx)
	if err != nil {
		return err
	}
	for _, segment := range segments {
		_ = s.logs.Write(ctx, &model.UserLog{
			UserID:      segment.UserID,
			Slug:        segment.Slug,
			Operation:   model.DeleteOp.String(),
			RequestTime: time.Now(),
		})
	}
	return nil
}
