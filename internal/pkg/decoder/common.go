package decoder

import "regexp"

func DecodeString(parser *regexp.Regexp, matchLabel string, matches []string) string {
	return matches[parser.SubexpIndex(matchLabel)]
}
