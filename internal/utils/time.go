package utils

import (
	"log"
	"time"
)

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("Asia/Makassar")
	if err != nil {
		log.Println("Error loading location:", err)
		return
	}
}

func IsAfterUnix(unixTime int64) bool {
	deadline := time.UnixMilli(unixTime).In(location)
	now := time.Now().In(location)
	return now.After(deadline)
}
