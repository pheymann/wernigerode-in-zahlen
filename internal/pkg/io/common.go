package io

import "os"

func WriteFile(filepath string, filename string, content string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		os.MkdirAll(filepath, 0700)
	}

	file, err := os.Create(filepath + filename)
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
