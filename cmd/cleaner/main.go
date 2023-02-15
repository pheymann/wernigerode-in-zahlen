package main

import (
	"os"

	"wernigode-in-zahlen.de/internal/cmd/cleaner"
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

	cleaner.CleanUp(filename, file, debug)
}

// func NewMetadataParser() MetadataParser {
// 	return MetadataParser{
// 		regexParser: []*regexp.Regexp{
// 			regexp.MustCompile(
// 				fmt.Sprintf(
// 					"^Dezernat/( )+Fachbereich (?P<department>\\d+),(?P<department_name>[ %s]+),+verantwortlich: (?P<accountable>[ %s]+)",
// 					rxGermanLetter,
// 					rxGermanLetter,
// 				),
// 			),
// 		},
// 	}
// }
