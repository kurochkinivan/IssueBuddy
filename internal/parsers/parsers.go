package parsers

import "strings"

func ParseLine(line string) []string {
	line = strings.ReplaceAll(line, " ", "")
	line = strings.TrimSpace(line)
	if line == "" {
		return []string{}
	}
	return strings.Split(line, ",")
}
