package decoder

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	RxGermanLetter            = "\\w\u00c4\u00e4\u00d6\u00f6\u00dc\u00fc\u00df"
	RxGermanPlusSpecialLetter = RxGermanLetter + `\(\);:\-/\.\d§&€><%`
)

func DecodeString(decoder *regexp.Regexp, matchLabel string, matches []string) string {
	return matches[decoder.SubexpIndex(matchLabel)]
}

func DecodeBudget(parser *regexp.Regexp, matchLabel string, matches []string) float64 {
	// 123.456,78 -> 123456.78
	strNumber := strings.ReplaceAll(
		strings.ReplaceAll(
			matches[parser.SubexpIndex(matchLabel)],
			".",
			"",
		),
		",",
		".",
	)

	i, err := strconv.ParseFloat(strNumber, 64)
	if err != nil {
		panic(err)
	}
	return i
}
