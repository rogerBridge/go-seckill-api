package common

import "log"

func Errlog(err error, msg string) {
	if err!=nil {
		log.Fatalf("%s: %s", err, msg)
	}
}