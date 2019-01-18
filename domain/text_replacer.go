package domain

import "bytes"

var (
	replaceFrom = "?"
	replaceTo = "!"
)

func ReplaceBytes(text []byte) []byte {
	return bytes.Replace(text, []byte(replaceFrom), []byte(replaceTo), -1)
}
