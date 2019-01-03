package main

import "strings"

// WithColors formats provided string with term colors.
func WithColors(str string) string {
	str = strings.Replace(str, "|x|", "\x1b[0m", -1)
	str = strings.Replace(str, "|w>", "\x1b[37m", -1)
	str = strings.Replace(str, "|_>", "\x1b[90m", -1)
	str = strings.Replace(str, "|r>", "\x1b[31m", -1)
	str = strings.Replace(str, "|y>", "\x1b[33m", -1)
	str = strings.Replace(str, "|g>", "\x1b[32m", -1)
	str = strings.Replace(str, "|b>", "\x1b[34m", -1)
	str = strings.Replace(str, "|m>", "\x1b[35m", -1)
	str = strings.Replace(str, "|c>", "\x1b[36m", -1)
	str = strings.Replace(str, "|B>", "\033[1m", -1)
	str = strings.Replace(str, "|N>", "\033[0m", -1)

	return str
}
