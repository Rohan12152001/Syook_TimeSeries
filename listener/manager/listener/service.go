package listener

import (
	"encoding/json"
	"fmt"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener/data"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener/db"
	utils2 "github.com/Rohan12152001/Syook_TimeSeries/listener/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type manager struct {
	db db.TimeSeriesDBManager
}

func New() ListenerManager {
	return manager{
		db: db.New(),
	}
}

var logger = logrus.New()

// DecryptAndEmit is mostly a Goroutine so no returns & ONLY panic ?  TODO: ASK THIS
func (m manager)DecryptAndEmit(enString, myListenerId string) (data.LiveData, error) {
	deStr, err := utils2.Decrypt(enString, utils2.SecretKey)
	if err != nil {
		logger.Error(err)
	}

	// 1. routine for validation
	var dataPayload data.LiveData
	bytes := []byte(deStr)

	err = json.Unmarshal(bytes, &dataPayload)
	if err != nil {
		logger.Error(err)
		return data.LiveData{}, err
	}

	// Validate here
	ok := utils2.ObjectValidation(dataPayload)
	if !ok{
		// Object discarded
		fmt.Println("Discarded!")
		return data.LiveData{}, DiscardedError
	}

	//2. Adding messages into DB
	newObjectId := uuid.New().String()

	dataBytes, err := json.Marshal(dataPayload)
	if err != nil {
		panic (err)
	}

	err = m.db.InsertObject(string(dataBytes), myListenerId, newObjectId)
	if err != nil {
		logger.Error(err)
		return data.LiveData{}, err
	}

	// 3. Return decrypted message
	return dataPayload, nil
}


