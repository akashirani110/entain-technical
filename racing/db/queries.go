package db

const (
	racesList = "list"
	raceById  = "raceById"
)

func getRaceQueries() map[string]string {
	return map[string]string{
		racesList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time,
				CASE WHEN (datetime('now', 'localtime') <= advertised_start_time) THEN 'OPEN' ELSE 'CLOSED' END AS status
			FROM races
		`,

		raceById: `
			SELECT
				id,
				meeting_id,
				name,
				number,
				visible,
				advertised_start_time,
				CASE WHEN (datetime('now', 'localtime') <= advertised_start_time) THEN 'OPEN' ELSE 'CLOSED' END AS status
			FROM races
			WHERE id = ?
			LIMIT 1
		`,
	}
}
