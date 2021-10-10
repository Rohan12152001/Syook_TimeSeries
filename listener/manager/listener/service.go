package listener

import (
	"encoding/json"
	"fmt"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener/data"
	utils2 "github.com/Rohan12152001/Syook_TimeSeries/listener/utils"
)

/*
Both steps need to be at service layer, and return the object here at ENDPOINT layer & use it
1. routine for validation
2. Adding messages into DB
3. sending decrypted messages to all UI's
*/

// DecryptAndEmit is mostly a Goroutine so no returns & ONLY panic ?  TODO: ASK THIS
func DecryptAndEmit(splitArray []string, liveChannel *chan data.LiveData) error {
	// Helper
	for _, enString := range splitArray {
		deStr, err := utils2.Decrypt(enString, utils2.SecretKey)
		if err != nil {
			panic(err)
		}

		// 1. routine for validation
		var dataPayload data.LiveData
		bytes := []byte(deStr)

		err = json.Unmarshal(bytes, &dataPayload)
		if err != nil {
			fmt.Println("Unmarshal problem...")
			return err
		}

		// Validate here
		ok := utils2.ObjectValidation(dataPayload)
		if !ok{
			// Object discarded
			fmt.Println("Discarded!")
			continue
		}

		//TODO: 2. Adding messages into DB
		// HERE

		// 3. Sending decrypted messages to all UI's
		*liveChannel <- dataPayload
	}
	close(*liveChannel)
	return nil
}


