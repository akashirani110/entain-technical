package db

import (
	"syreclabs.com/go/faker"
	"time"
)

func (s *sportsRepo) seed() error {
	statement, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS sports_events (id INTEGER PRIMARY KEY, name TEXT, advertised_start_time DATETIME, location TEXT, description TEXT, status TEXT, winner TEXT)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = s.db.Prepare(`INSERT OR IGNORE INTO sports_events(id, name, advertised_start_time, location, description) VALUES (?, ?, ?, ?, ?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				getRandomSports(),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 3)).Format(time.RFC3339),
				faker.Address().City(),
				faker.Lorem().Paragraph(3),
			)
		}
	}

	return err
}

func getRandomSports() string {
	sports := []string{
		"Football",
		"Basketball",
		"Tennis",
		"Golf",
		"Cricket",
		"Baseball",
		"Hockey",
		"Soccer",
		"Badminton",
		"Volleyball",
		"Ice Hockey",
		"Rugby",
		"Footy",
		"Polo",
	}
	return sports[faker.RandomInt(0, len(sports)-1)]
}
