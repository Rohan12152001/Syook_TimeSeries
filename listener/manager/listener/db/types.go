package db

import (
	_ "github.com/lib/pq"
	"os"
)

type TimeSeriesDBManager interface {
	InsertObject(dataPayload, myListenerId, newObjectId string) error
}


// handle getting from env gracefully
var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = 5432
	user     = "postgres"
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)
