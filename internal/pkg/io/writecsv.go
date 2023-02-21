package io

import (
	"encoding/csv"
	"os"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func WriteCSV(target model.TargetFile, content [][]string) {
	if _, err := os.Stat(target.Path); os.IsNotExist(err) {
		os.MkdirAll(target.Path, 0700)
	}

	file, err := os.Create(target.CanonicalName())
	if err != nil {
		panic(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	for _, record := range content {
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}
	writer.Flush()
}
