package server

import (
	"errors"
	"github.com/dustinmj/renotts/coms"
)

//AvailServs - map of services available
var AvailServs = map[string]Eng{
	"polly": Polly}

/*Eng tts engine interface
Query returns the sound file and http response object
SetDefs allows engine to write it's own defaults, these can then be used and will
be honored if they exist in the main config file from user
Caches return true if the Sf file returned from query will need to be played,
it will then be passed to the audio player (mp3);
return false if the engine handles playing the file internally
*/
type Eng interface {
	Query(*Rq) (Sf, Rsp)
	SetDefs()
	Caches() bool
}

//Serv - structure for implementing engine interface
type engine struct{}

func enGet(t string) (e Eng, er error) {
	for n, s := range AvailServs {
		if n == t {
			s.SetDefs()
			return s, nil
		}
	}
	return nil, errors.New(coms.Err["NoService"])
}
