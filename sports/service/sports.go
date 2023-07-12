package service

import (
	"golang.org/x/net/context"
	"sports/db"
	"sports/proto/sports"
)

type Sports interface {
	// ListEvents will return a collection of sports events.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)
}

// sportsService implements the Sports interface
type sportsService struct {
	sportsRepo db.SportsRepo
}

// NewSportsService instantiates and returns a new sportsService.
func NewSportsService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}

func (s *sportsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	sportsEvent, err := s.sportsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &sports.ListEventsResponse{Events: sportsEvent}, nil
}
