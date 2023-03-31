package metadata

import (
	"wernigerode-in-zahlen.de/internal/pkg/io"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Write(metadata string, target model.TargetFile) {
	target.Tpe = "json"

	io.WriteFile(target, metadata)
}
