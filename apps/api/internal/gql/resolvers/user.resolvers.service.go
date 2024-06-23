package resolvers

import (
	"regexp"
)

var userNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
