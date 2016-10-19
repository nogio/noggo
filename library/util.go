package library

import (
	"encoding/base64"
)

func encode64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
func decode64(value string) string {
	d, e := base64.StdEncoding.DecodeString(value)
	if e == nil {
		return string(d)
	}
	return value
}


