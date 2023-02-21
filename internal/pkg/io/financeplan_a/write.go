package filewriter

import (
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func WriteGroup(financePlans [][]string, target model.TargetFile) {
	target.Tpe = "csv"

	io.WriteCSV(target, financePlans)
}

func WriteUnit(financePlans map[string][][]string, target model.TargetFile) {
	for costCenterUnit, financePlans := range financePlans {
		if len(financePlans) == 0 {
			continue
		}

		targetCpy := target
		targetCpy.Tpe = "csv"
		targetCpy.Path = target.Path + costCenterUnit + "/"

		io.WriteCSV(targetCpy, financePlans)
	}
}
