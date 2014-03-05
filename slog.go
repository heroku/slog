// Some Strcutred Logging Help

package slog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	Escapers = regexp.MustCompile(`["'\n\s/=:,]`)
)

type Context map[string]interface{}

// Does the hard work of converting the context
// to an alphabetically key sorted string
func (c Context) String() string {
	var sv string
	parts := make([]string, 0, len(c))

	for k, v := range c {
		switch v.(type) {
		case time.Time: // Format times the way we want them
			sv = v.(time.Time).Format(time.RFC3339Nano)
		case time.Duration:
			t := v.(time.Duration)
			switch {
			case t < time.Microsecond:
				sv = fmt.Sprintf("%.9fs", v.(time.Duration).Seconds())
			case t < time.Millisecond:
				sv = fmt.Sprintf("%.6fs", v.(time.Duration).Seconds())
			default:
				sv = fmt.Sprintf("%.3fs", v.(time.Duration).Seconds())
			}
		case int:
			sv = fmt.Sprintf("%d", v.(int))
		case string:
			sv = fmt.Sprintf("%s", v.(string))
		case error:
			sv = fmt.Sprintf("%s", v.(error))
		default: // Let Go figure out the representation
			sv = fmt.Sprintf("%v", v)
		}

		// If there are any spaces characters then need to quote the value
		if Escapers.MatchString(sv) {
			sv = fmt.Sprintf("%q", sv)
		}

		// Assemble the final part and append it to the array
		parts = append(parts, fmt.Sprintf("%s=%s", k, sv))
	}
	sort.Strings(parts)
	return strings.Join(parts, " ")
}

// Pushes count#what=value onto the context
// if count#what already exists in the context
// value is added to it
func (c Context) Count(what string, value int) {
	what = fmt.Sprintf("count#%s", what)
	ov, ok := c[what]
	if ok {
		value = ov.(int) + value
	}
	c[what] = value
}

// Pushes measure#what=value onto the context
func (c Context) Measure(what string, value interface{}) {
	c[fmt.Sprintf("measure#%s", what)] = value
}

// Pushes sample#what=value onto the context
func (c Context) Sample(what string, value interface{}) {
	c[fmt.Sprintf("sample#%s", what)] = value
}

// Pushes unique#what=value onto the context
func (c Context) Unique(what string, value interface{}) {
	c[fmt.Sprintf("unique#%s", what)] = value
}

// Method wrapper
func (c Context) Add(what string, value interface{}) {
	c[what] = value
}