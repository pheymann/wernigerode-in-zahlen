package targetfile

import (
	"fmt"
	"os"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

var (
	fileRegex = regexp.MustCompile(
		`^assets/data/raw/(?P<path>(\d+/)+)(?P<file_name>\w+)\.(?P<file_type>\w+)`,
	)
)

func Decode(source *os.File) model.TargetFile {
	matches := fileRegex.FindStringSubmatch(source.Name())
	path := matches[fileRegex.SubexpIndex("path")]
	fileName := matches[fileRegex.SubexpIndex("file_name")]

	return model.TargetFile{
		Path: fmt.Sprintf("assets/data/processed/%s", path),
		Name: fileName,
	}
}
