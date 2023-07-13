package db

import (
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"sports/proto/sports"
	"sync"
	"syreclabs.com/go/faker"
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
	// select sportsList query lists all the sports events
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
		var statusString string
		var winner sql.NullString

		if err := rows.Scan(&event.Id, &event.Name, &advertisedStart, &event.Location, &event.Description, &statusString, &winner); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		// convert the status string to the corresponding enum value
		statusValue := sports.EventStatus_value[statusString]

		statusEnum := sports.EventStatus(statusValue)

		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts
		// update the field
		event.Status = &statusEnum
		// check if the winner value is valid (not null)
		if winner.Valid {
			event.Winner = &winner.String
		} else {
			// check if the status is COMPLETED
			if statusEnum == sports.EventStatus_COMPLETED {
				// set the winner of the completed sports event
				winner.String, err = generateRandomWinner()
				if err != nil {
					return nil, err
				}
				event.Winner = &winner.String
			} else {
				// else set the winner as nil
				event.Winner = nil
			}
		}

		events = append(events, &event)
	}

	return events, nil
}

// generates random team names for the winner
func generateRandomWinner() (string, error) {
	return faker.Team().Name(), nil
}
