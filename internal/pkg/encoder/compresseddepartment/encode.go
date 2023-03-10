package compresseddepartment

import (
	"encoding/json"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func Encode(compressed model.CompressedDepartment) string {
	bytes, err := json.MarshalIndent(compressed, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}