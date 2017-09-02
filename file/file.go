package file

import (
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/tmplt"
	"os"
	"path/filepath"
)

//SilencePath silence file full path
var SilencePath string

//SystemdPath service File
var SystemdPath string

const (
	//Smp3 Silence mp3 file
	smp3 = "Silence.mp3"
	//Systemd service file name
	systemd = "renotts.service"
)

//Setup setup file requirements
func Setup(cfg config.Cfg) {
	// attempt to ensure silence.mp3 is expanded to cachepath
	SilencePath = filepath.Join(cfg.Cache(), smp3)
	putSilence(SilencePath, cfg.Cache())
	// ensure init.d utilities have been created
	SystemdPath = filepath.Join(cfg.Path(), systemd)
	putInitD(SystemdPath, cfg)
}

// expand silence file
func putSilence(silence string, cache string) {
	// expand silence.mp3 to file if necessary
	if !chkFile(silence) {
		if err := RestoreAsset(cache, smp3); err != nil {
			coms.Msg("Unable to restore Silence.mp3 to cachepath!")
		}
	}
}

func putInitD(path string, cfg config.Cfg) error {
	// attempt to create
	var f *os.File
	for i := 0; i < 2; i++ {
		var err error
		f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			// attempt to create directory on first err
			if i == 0 {
				err = os.MkdirAll(filepath.Dir(path), 0755)
				if err != nil {
					return err
				}
				continue
			} else {
				// don't loop on additional errors
				return err
			}
		}
		defer f.Close()
		break
	}
	// erase contents
	f.Truncate(0)
	// create new init.d
	data := tmplt.SysD{
		User:    cfg.User(),
		AppName: cfg.AppName(),
		AppPath: cfg.AppPath()}
	return tmplt.ParseF(f, tmplt.SystemdFl, data)
}

func chkFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
