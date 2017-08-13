package coms

// Err - contains error messages
var Err = map[string]string{
	"InvalidJSONInBody": "cannot read JSON data in POST body",
	"ErrorReadingBody":  "could not retrieve POST body",
	"NoService":         "service does not exist",
	"InvalidPath":       "invalid path specified"}
