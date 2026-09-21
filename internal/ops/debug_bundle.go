package ops

import "strings"

var secretMarkers = []string{"TOKEN", "PASSWORD", "SECRET", "KEY", "CREDENTIAL"}

func RedactEnv(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		name := line
		if i := strings.IndexByte(line, '='); i >= 0 {
			name = line[:i]
		}
		if secretName(name) {
			out = append(out, name+"=<redacted>")
			continue
		}
		out = append(out, line)
	}
	return out
}

func secretName(name string) bool {
	upper := strings.ToUpper(name)
	for _, marker := range secretMarkers {
		if strings.Contains(upper, marker) {
			return true
		}
	}
	return false
}
