package metadata

import (
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func Write(metadata string, target model.TargetFile) {
	target.Tpe = "json"

	io.Write(target, metadata)
}
