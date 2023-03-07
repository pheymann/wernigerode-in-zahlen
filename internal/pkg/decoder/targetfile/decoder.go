package targetfile

import (
	"fmt"
	"os"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

var (
	fileRegex = regexp.MustCompile(
		`^assets/data/(raw)?(processed)?/(?P<path>(\d+/)+)(?P<file_name>\w+)\.(?P<file_type>\w+)`,
	)
)

func Decode(source *os.File, tpe string) model.TargetFile {
	matches := fileRegex.FindStringSubmatch(source.Name())
	path := matches[fileRegex.SubexpIndex("path")]
	fileName := matches[fileRegex.SubexpIndex("file_name")]

	return model.TargetFile{
		Path: fmt.Sprintf("assets/%s/%s", tpe, path),
		Name: fileName,
	}
}
