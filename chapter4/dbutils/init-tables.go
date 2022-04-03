package dbutils

import "database/sql"

func Initialize(db *sql.DB) error {
	for _, table := range []string{train, station, schedule} {
		statement, err := db.Prepare(table)
		if err != nil {
			return err
		}
		_, err = statement.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
