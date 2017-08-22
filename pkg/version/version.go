package version

import (
	"strconv"
	"time"
)

// Version is the version string from the application. This variable is set
// during the build process.
var Version = "UNKNOWN"

// BuildDate is the build date for the application. This variable is set during
// the build process.
var BuildDate = "UNKNOWN"

// ParseBuildDate parses the unix timestamp value set in the BUILDDATE
// variable. It returns either a time or an error.
func ParseBuildDate() (time.Time, error) {
	d, err := strconv.ParseInt(BuildDate, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(d, 0), nil
}
