package utils

import "log"

//HandleError handles error
func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
