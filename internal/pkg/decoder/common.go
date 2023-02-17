package decoder

import "regexp"

const (
	RxGermanLetter = "\\w\u00c4\u00e4\u00d6\u00f6\u00dc\u00fc\u00df"
)

func DecodeString(decoder *regexp.Regexp, matchLabel string, matches []string) string {
	return matches[decoder.SubexpIndex(matchLabel)]
}
