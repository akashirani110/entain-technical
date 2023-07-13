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
				description,
				CASE
        			WHEN (datetime('now', 'localtime') < advertised_start_time) THEN 'UPCOMING'
        			WHEN (datetime('now', 'localtime') >= advertised_start_time) THEN
            			CASE
                			WHEN (datetime('now', 'localtime') > advertised_start_time) THEN 'COMPLETED'
                			ELSE 'IN_PROGRESS'
            			END
        			ELSE 'STATUS_UNSPECIFIED'
    			END AS status,
			    winner
			FROM sports_events
		`,
	}
}
