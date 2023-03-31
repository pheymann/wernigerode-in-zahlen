package metadatacleaner

import (
	"bufio"
	"fmt"
	"os"

	decodeMeta "wernigerode-in-zahlen.de/internal/pkg/decoder/metadata"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Cleanup(metadataFile *os.File) model.Metadata {
	metadataScanner := bufio.NewScanner(metadataFile)
	metadataLines := []string{}

	for metadataScanner.Scan() {
		metadataLines = append(metadataLines, metadataScanner.Text())
	}

	metadataDecoder := decodeMeta.NewMetadataDecoder()

	defer func() {
		if r := recover(); r != nil {
			metadataDecoder.Debug()
			fmt.Printf("\n%+v\n", r)
			os.Exit(1)
		}
	}()

	metadata := metadataDecoder.DecodeFromCSV(metadataLines)

	return metadata
}
