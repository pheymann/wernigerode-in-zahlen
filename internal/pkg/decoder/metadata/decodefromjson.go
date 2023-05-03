package metadata

import (
	"encoding/json"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func DecodeFromJSON(metadataJSON string) model.Metadata {
	var metadata = model.Metadata{}
	json.Unmarshal([]byte(metadataJSON), &metadata)

	return metadata
}
