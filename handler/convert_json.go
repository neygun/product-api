package handler

import (
	"encoding/json"
	"log"
)

func ToJsonString(i interface{}) string {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		log.Fatalf("Error marshaling struct: %v", err)
	}
	return string(jsonBytes)
}

func ToStruct(s string) interface{} {
	var i interface{}
	err := json.Unmarshal([]byte(s), &i)
	if err != nil {
		log.Fatalf("Error unmarshaling: %v", err)
	}
	return i
}
