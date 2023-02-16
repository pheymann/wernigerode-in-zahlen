package main

import (
	"bufio"
	"fmt"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/cleaner"
	"wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
)

var (
	debug = false
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	metadataDecoder := metadata.NewMetadataDecoder()
	fmt.Println(metadataDecoder.Decode(lines))

	cleaner.CleanUp(filename, file, debug)
}
