package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type manager struct {
	db *sqlx.DB
}

var logger = logrus.New()

func New() TimeSeriesDBManager {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		logger.Error(err)
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

