package main

import (
	"fmt"
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/file"
	"github.com/dustinmj/renotts/server"
	"os"
	"os/user"
)

func main() {
	// A simple warning about running
	// RenoTTS as elevated user
	chkElevated()
	putSilence()
	server.Create()
}

func chkElevated() {
	warn := "!!!!! RenoTTS appears to be running as root user. This is not recommended. RenoTTS should not be run as root user. !!!!!"
	user, err := user.Current()
	if err != nil {
		return
	}
	if user.Uid == "0" || user.Name == "root" || user.Name == "su" || user.Name == "sudoer" {
		fmt.Println("")
		fmt.Println(warn)
		fmt.Println("")
	}
}

func putSilence() {
	// lint shadow warning
	var err error
	// make sure we have our silence file.
	config.SilenceFile, err = file.FullPath(config.Smp3)
	if err != nil {
		// handled in config
		coms.Msg("Unable to access cache.")
	}
	// expand silence.mp3 to file if necessary
	if _, err := os.Stat(config.SilenceFile); os.IsNotExist(err) {
		if err := file.RestoreAsset(config.Val(config.CACHEPATH), config.Smp3); err != nil {
			coms.Msg("Unable to restore Silence.mp3 to cachepath!")
		}
	}
}
