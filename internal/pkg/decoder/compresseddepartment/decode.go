package compresseddepartment

import (
	"encoding/json"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Decode(compressedJSON string) model.CompressedDepartment {
	var compressedDepartment model.CompressedDepartment
	json.Unmarshal([]byte(compressedJSON), &compressedDepartment)

	return compressedDepartment
}
