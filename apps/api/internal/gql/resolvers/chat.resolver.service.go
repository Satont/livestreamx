package resolvers

import (
	"regexp"
)

var badSymbolsRegexp = regexp.MustCompile("[^\\p{Cyrillic}\\p{Latin}\\p{Common}\\w\\s]+")
