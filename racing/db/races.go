package db

import (
	"database/sql"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, sortBy *racing.ListRacesSortBy, order *racing.ListRacesOrder) ([]*racing.Race, error)

	// Get will return a race by its id.
	Get(id int64) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

type SortBy int64
type Order int64

// enum for sort by column
const (
	Id SortBy = iota + 1
	MeetingId
	Name
	Number
	AdvertisedStartTime
)

// enum for order
const (
	Ascending Order = iota + 1
	Descending
)

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, sortBy *racing.ListRacesSortBy, order *racing.ListRacesOrder) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	query = r.applySortByAndOrder(query, sortBy, order)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}
	// apply visibility filter if passed true
	if filter.OnlyVisible {
		args = append(args, true)
		clauses = append(clauses, "visible="+strconv.FormatBool(filter.OnlyVisible))
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

// returns a string value for the enum value passed
func (sb SortBy) sortByString() string {
	return [...]string{"id", "meeting_id", "name", "number", "advertised_start_time"}[sb-1]
}

// apply order by and order to the query if available
// by default applies order by advertised_start_time in ascending order
func (r *racesRepo) applySortByAndOrder(query string, sortBy *racing.ListRacesSortBy, order *racing.ListRacesOrder) string {
	var sb SortBy
	if sortBy == nil {
		return query
	}

	// check if key is present
	if sortBy.Column != 0 {
		sb = SortBy(sortBy.Column)
		// order by given column key
		query += " ORDER BY " + sb.sortByString()
	} else {
		// sets advertised_start_time as default
		sb = SortBy(5)
		query += " ORDER BY " + sb.sortByString()
	}

	// check if order is Descending
	if order.Order == 2 {
		query += " DESC"
	} else {
		// sets order ascending as default
		query += " ASC"
	}

	return query
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &race.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}

	return races, nil
}

// Get takes id as input and returns a single record of Race if found
func (r *racesRepo) Get(id int64) (*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[raceById]
	args = append(args, id)

	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	races, _ := r.scanRaces(rows)

	// If any error occurs while iterating over the rows, return err
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If no race is found, return nil
	if len(races) == 0 {
		return nil, nil
	}

	// If a race is found, return it
	return races[0], nil
}
