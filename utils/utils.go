package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

//HandleError handles error
func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

//ToBytes encode i to slice of bytes
func ToBytes(i interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	HandleError(encoder.Encode(i))
	return buffer.Bytes()
}

//FromBytes decode slice of bytes to i
func FromBytes(i interface{}, data []byte) {
	HandleError(gob.NewDecoder(bytes.NewReader(data)).Decode(i))
}

//Hash return hash result of i
func Hash(i interface{}) string {
	toString := fmt.Sprintf("%v", i)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(toString)))
}

func Splitter(s, sep string, i int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}

func ToJSON(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleError(err)
	return r
}
