package server

import (
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
	"io"
	"os"
)

const silence = "Silence.mp3"

// created on init
var silenceFile string

//Aq audio query
type Aq struct {
	Txt    string
	Typ    string
	Buffer *io.ReadCloser
	Chars  int64
}

type sPlayer interface {
	play(*Sf) error
}

//Serv - structure for implementing engine interface
type player struct{}

func init() {
	// lint shadow warning
	var err error
	// make sure we have our silence file.
	silenceFile, err = FullPath(silence)
	if err != nil {
		// handled in config
		coms.Msg("Unable to access cache.")
	}
	if _, err := os.Stat(silenceFile); os.IsNotExist(err) {
		if err := RestoreAsset(config.Val("cachepath"), silence); err != nil {
			coms.Msg("Unable to restore Silence.mp3 to cachepath!")
		}
	}
}
