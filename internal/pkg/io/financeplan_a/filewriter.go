package filewriter

import (
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func WriteGroup(financePlans string, target model.TargetFile) {
	target.Tpe = "csv"

	io.Write(target, financePlans)
}

func WriteUnit(financePlans map[string]string, target model.TargetFile) {
	for _, financePlans := range financePlans {
		if len(financePlans) == 0 {
			continue
		}

		target.Tpe = "csv"

		io.Write(target, financePlans)
	}
}
