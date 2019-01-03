package main

import "fmt"

// Error is main error type for this util
type Error struct {
	ID string
}

// NewCoreError - create new error from given error-id
func NewCoreError(id string) Error {
	return Error{
		ID: id,
	}
}

// PrintError prints error
func PrintError(id string) {
	var errMsg string

	switch id {
	case "io":
		errMsg = "IO Error"
	case "json":
		errMsg = "Error with [de]serializing json"
	case "msgpack-encoding":
		errMsg = "Error with encoding MsgPack"
	case "msgpack-decoding":
		errMsg = "Error with decoding MsgPack"
	case "crypto":
		errMsg = "Encryption/Decryption error"
	case "key-len":
		errMsg = "Wrong key length error"
	case "con":
		errMsg = "Communication error"
	case "internal":
		errMsg = "Mysterious internal error ԅ(≖‿≖ԅ)"
	case "not-found":
		errMsg = "Cannot find that... ¯\\_(ツ)_/¯"
	case "incorrenct-config":
		errMsg = "Incorrect config (╯°□°）╯︵ ┻━┻"
	case "incorrenct-request":
		errMsg = "Incorrect request"
	case "duplicate":
		errMsg = "Duplicate"
	case "invalid-key":
		errMsg = "Invalid key ( ° ͜ʖ͡°)╭∩╮"
	case "not-enough-keys":
		errMsg = "Not enough keys ( ° ͜ʖ͡°)╭∩╮"
	case "unknown":
		errMsg = "Unknown error ( ﾉ^ω^)ﾉﾟ"
	default:
		errMsg = fmt.Sprintf("Another error: '%s'", id)
	}

	fmt.Printf(WithColors("|r>✘ %s|x|\n"), errMsg)
}

func (e Error) Error() string {
	return e.ID
}
