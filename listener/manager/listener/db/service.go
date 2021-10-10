package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type manager struct {
	db *sqlx.DB
}

func New() TimeSeriesDBManager {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return manager{
		db: db,
	}
}

func (m manager) InsertObject(dataPayload, myListenerId, newObjectId string) error {
	query := "Insert into TimeSeries (id, listenerId, data) values($1, $2, $3);"

	_, err := m.db.Exec(query, newObjectId, myListenerId, dataPayload)
	if err != nil {
		return err
	}

	return nil
}

