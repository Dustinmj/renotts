package player

import (
	"github.com/dustinmj/renotts/config"
)

//SPlayer interface for players
type SPlayer interface {
	Play(string, bool, bool) error
	Busy() bool
	Queue(string, bool, bool) error
}

type playerQueueFile struct {
	Path   string
	Before bool
	After  bool
}

//GetPlayer returns the proper player based on config settings
func GetPlayer() SPlayer {
	if config.Exists(config.EXECPLAYER) && len(config.Val(config.EXECPLAYER)) > 0 {
		return &execPlayer
	}
	return &mpgPlayer
}
