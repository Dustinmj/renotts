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
func GetPlayer(cfg config.Cfg) SPlayer {
	if cfg.Exists(config.EXECPLAYER) && len(cfg.Val(config.EXECPLAYER)) > 0 {
		return &eplayer{
			player: cfg.Val(config.EXECPLAYER)}
	}
	if cfg.Exists(config.USERBUFFERSIZE) {
		return &mplayer{
			userBufferSize: cfg.Intval(config.USERBUFFERSIZE)}
	}
	return nil
}
