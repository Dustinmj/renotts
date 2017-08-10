package aud

import (
	"github.com/dustinmj/renotts/com"
	"github.com/dustinmj/renotts/config"
	"os/exec"
)

//PlayExec execute a play command through OS
func PlayExec(sF com.Sf) error { // TODO
	com.Msg("Playing file: " + sF.Path)
	cmd := exec.Command(config.Val("execCmd"), sF.Path)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
