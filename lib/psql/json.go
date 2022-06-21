package psql

import (
	"encoding/json"
)

type JSON map[string]string

func (j JSON) Add(key string, value string) JSON {
	j[key] = value
	return j
}

func (j JSON) Marshal() []byte {
	js, _ := json.Marshal(j)
	return js
}

func (j JSON) GetValue(key string) string {

	if len(key) == 0 {
		return ""
	}

	value, ok := j[key]
	if !ok {
		return ""
	}

	return value
}

func (j JSON) String() string {
	return string(j.Marshal())
}
