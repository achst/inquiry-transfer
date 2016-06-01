package json

import (
	"encoding/json"
)

var j interface{}

func (j *json) StrToJson(s string) {
	err := json.Unmarshal(s, &j)
	if err != nil {
		panic(err)
	}
	println(j)
}
