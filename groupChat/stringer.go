package main

import "encoding/json"

func toString(o interface{}) string {
	s, _ := json.MarshalIndent(o, "", "\t")
	return string(s)
}
