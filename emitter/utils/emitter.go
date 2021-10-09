package utils

import (
	"4d63.com/optional"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type messageStruct struct {
	Name string `json:"name"`
	Origin string `json:"origin"`
	Destination string `json:"destination"`
	Secret_key string `json:"secret_key"`
}

func getRandomInt(upperLimit int, exclude optional.Int64) int64{
	if !exclude.IsPresent(){
		return rand.Int63n(int64(upperLimit))
	}

	excludeValue, _ := exclude.Get()
	for {
		r := rand.Int63n(int64(upperLimit))
		if r != excludeValue {
			return r
		}
	}
}

// Form random object
func formRandomObject() messageStruct{
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	randomName := AllData.Names[getRandomInt(len(AllData.Names), nil)]

	originIndex := getRandomInt(len(AllData.Cities), nil)
	randomOrigin := AllData.Cities[originIndex]

	randomDestination := AllData.Cities[getRandomInt(len(AllData.Cities), optional.OfInt64(originIndex))]

	message := messageStruct{
		Name: randomName,
		Origin: randomOrigin,
		Destination: randomDestination,
	}

	fmt.Println(message)
	return message
}


// Add secret-hash & Generate string func
func formRandomObjectString() string{
	randomObject := formRandomObject()

	// Form secret hash & Add in struct
	secretHash := ""

	byteRandomObject, err := json.Marshal(randomObject)
	if err != nil {
		fmt.Errorf("Error!")
		return ""
	}

	hasher := sha1.New()
	hasher.Write(byteRandomObject)
	secretHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	randomObject.Secret_key = secretHash
	fmt.Println(randomObject)

	// Form final string
	bytes, _ := json.Marshal(randomObject)
	fmt.Println(bytes)

	return string(bytes)
}

// Return the long encrypted string
func FormFinalString() string{
	finalString := ""

	rand.Seed(time.Now().Unix())
	for count:=0;count<rand.Intn(499 - 49) + 59;count++{
		tempString := formRandomObjectString()
		if count==0{
			finalString += tempString
		}else{
			finalString += "|" + tempString
		}
	}

	return finalString
}