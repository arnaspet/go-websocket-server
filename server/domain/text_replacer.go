package domain

import "bytes"

const (
	replaceFrom = "?"
	replaceTo = "!"
)

func ReplaceBytes(text []byte) []byte {
	return bytes.Replace(text, []byte(replaceFrom), []byte(replaceTo), -1)
}
