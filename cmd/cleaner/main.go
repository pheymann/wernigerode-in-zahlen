package main

import (
	"bufio"
	"fmt"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/cleaner"
	decodeMeta "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	encodeMeta "wernigode-in-zahlen.de/internal/pkg/encoder/metadata"
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

	metadataDecoder := decodeMeta.NewMetadataDecoder()
	fmt.Printf("%+v\n", string(encodeMeta.Encode(metadataDecoder.Decode(lines))))

	cleaner.CleanUp(filename, file, debug)
}
