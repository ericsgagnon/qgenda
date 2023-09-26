package qgenda

import (
	"encoding/json"
	"log"
	"os"
)

func GetFromFile[P *[]T, T any](filename string, dst P) error {

	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	t := []T{}
	if err := json.Unmarshal(b, &t); err != nil {
		log.Println(err)
	}
	*dst = t
	return nil

}

func toPointer[P *V, V any](value V) P {
	return &value
}

func toValue[P *V, V any](pointer P) V {
	return *pointer
}
