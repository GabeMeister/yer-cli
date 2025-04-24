package utils

import "encoding/json"

func MustJSONStringify(v interface{}) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
