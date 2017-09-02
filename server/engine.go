package server

import (
	"errors"
	"github.com/dustinmj/renotts/config"
)

//AvailServs map of services available
var AvailServs = map[string]eng{
	"polly": Polly}

//Eng tts engine interface
type eng interface {
	Query(*request, config.Cfg) (*string, error)
	SetDefs()
	Caches() bool
}

//Serv - structure for implementing engine interface
type engine struct{}

func enGet(t string) (e eng, er error) {
	for n, s := range AvailServs {
		if n == t {
			s.SetDefs()
			return s, nil
		}
	}
	return nil, errors.New(errBadServce)
}
