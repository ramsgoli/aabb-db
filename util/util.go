package util

import "strings"

func Trim(b []byte) string {
	return strings.TrimRight(string(b), "\x00")
}
