package db

import (
	_ "github.com/lib/pq"
)

type TimeSeriesDBManager interface {
	InsertObject(dataPayload, myListenerId, newObjectId string) error
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "Syook_Timeseries"
)
