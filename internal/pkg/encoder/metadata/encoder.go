package metadata

import (
	"encoding/json"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Encode(metadata model.Metadata) string {
	bytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
