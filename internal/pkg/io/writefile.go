package io

import (
	"os"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func WriteFile(target model.TargetFile, content string) {
	if _, err := os.Stat(target.Path); os.IsNotExist(err) {
		os.MkdirAll(target.Path, 0700)
	}

	file, err := os.Create(target.CanonicalName())
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
	file.Sync()
}
