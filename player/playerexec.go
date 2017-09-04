package player

import (
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/file"
	"os/exec"
)

//structure for implementing engine interface
type eplayer struct{}

var execPlayer = eplayer{}
var execPlayerQueue = []playerQueueFile{}
var execIsPlaying bool

func (execPlayer *eplayer) Play(path string, padB bool, padA bool, player string) error {
	execIsPlaying = true
	coms.Msg("Playing file", "via", player+":", path)
	// before silence
	if padB {
		if err := execPlayer.playSilence(player); err != nil {
			execIsPlaying = false
			return err
		}
	}
	if err := execPlayer.execCommand(player, path); err != nil {
		execIsPlaying = false
		return err
	}
	// after silence
	if padA {
		if err := execPlayer.playSilence(player); err != nil {
			execIsPlaying = false
			return err
		}
	}
	execIsPlaying = false
	coms.Msg("Completed Playing")
	if len(execPlayerQueue) > 0 {
		next := execPlayerQueue[0]
		execPlayerQueue = execPlayerQueue[1:]
		execPlayer.Play(next.Path, next.Before, next.After, player)
	}
	return nil
}

func (execPlayer *eplayer) Busy() bool {
	return execIsPlaying
}

func (execPlayer *eplayer) Queue(path string, before bool, after bool) error {
	execPlayerQueue = append(execPlayerQueue,
		playerQueueFile{
			Path:   path,
			Before: before,
			After:  after})
	return nil
}

func (execPlayer *eplayer) execCommand(cmd string, path string) error {
	com := exec.Command(cmd, path)
	err := com.Start()
	if err != nil {
		return err
	}
	err = com.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (execPlayer *eplayer) playSilence(player string) error {
	if err := execPlayer.execCommand(player, file.SilencePath); err != nil {
		return err
	}
	return nil
}
