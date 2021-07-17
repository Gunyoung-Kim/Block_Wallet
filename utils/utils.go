package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

//HandleError handles error
func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	HandleError(encoder.Encode(i))
	return buffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	HandleError(gob.NewDecoder(bytes.NewReader(data)).Decode(i))
}
