package resolvers

import (
	"regexp"
)

var userNameRegexp = regexp.MustCompile(`^[a-z0-9_]+$`)
var userDisplayNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
