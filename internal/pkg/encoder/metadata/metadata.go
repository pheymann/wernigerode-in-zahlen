package metadata

import (
	"encoding/json"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func Encode(metadata model.Metadata) []byte {
	bytes, err := json.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	return bytes
}
