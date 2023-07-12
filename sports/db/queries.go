package db

const (
	sportsList = "sportsList"
)

func getEventsQueries() map[string]string {
	return map[string]string{
		sportsList: `
			SELECT
				id,
				name,
				advertised_start_time,
				location,
				description
			FROM sports_events
		`,
	}
}
