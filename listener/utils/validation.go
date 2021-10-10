package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener/data"
)

func ObjectValidation(dataPayload data.LiveData) bool{
	validatePayload := data.LiveData{
		Name: dataPayload.Name,
		Origin: dataPayload.Origin,
		Destination: dataPayload.Destination}

	// Form secret hash
	secretHash := ""

	ObjectBytes, err := json.Marshal(validatePayload)
	if err != nil {
		fmt.Errorf("Error!")
		return false
	}

	hasher := sha1.New()
	hasher.Write(ObjectBytes)
	secretHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Check
	if secretHash==dataPayload.SecretKey{
		return true
	}
	return false
}


