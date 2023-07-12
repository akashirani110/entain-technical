package db

import (
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"sports/proto/sports"
	"sync"
	"time"
)

// SportsRepo provides repository access to sports events.
type SportsRepo interface {
	Init() error
	List(filter *sports.ListEventsFilter) ([]*sports.SportsEvent, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sports repository dummy data.
func (s *sportsRepo) Init() error {
	var err error

	s.init.Do(func() {
		err = s.seed()
	})

	return err
}

// List queries the database and returns the records from scanEvents
func (s *sportsRepo) List(filter *sports.ListEventsFilter) ([]*sports.SportsEvent, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventsQueries()[sportsList]

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanEvents(rows)
}

// scanEvents scans the db rows and returns list of sports events
func (s *sportsRepo) scanEvents(rows *sql.Rows) ([]*sports.SportsEvent, error) {
	var events []*sports.SportsEvent

	for rows.Next() {
		var event sports.SportsEvent
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Name, &advertisedStart, &event.Location, &event.Description); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts

		events = append(events, &event)
	}

	return events, nil
}
