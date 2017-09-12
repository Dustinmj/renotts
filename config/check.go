package config

import (
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
)

type chk struct{}

//ConfigChk - check configurations
var ConfigChk = chk{}

func (ConfigChk chk) All() []string {
	conf := Get()
	var s []string
	s = append(s, ConfigChk.Config(conf)...)
	s = append(s, ConfigChk.Cache(conf))
	s = append(s, ConfigChk.Player(conf)...)
	s = append(s, ConfigChk.Amazon(conf)...)
	return s
}

func (ConfigChk chk) Config(conf Cfg) []string {
	var s []string
	// check config file to make sure it exists
	cfg := conf.File()
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		s = append(s, "missing RenoTTS configuration file: "+cfg)
	}
	s = append(s, "found RenoTTS configuration file: "+cfg)
	// check config values
	vals := viper.AllKeys()
	for _, k := range vals {
		s = append(s, "---> current configuration ---> "+k+": "+conf.Val(k))
	}
	return s
}

func (ConfigChk chk) Cache(conf Cfg) string {
	// check config file to make sure it exists
	test := filepath.Join(conf.Val(CACHEPATH), "test.txt")
	if _, err := os.OpenFile(test, os.O_RDWR|os.O_CREATE, 0755); err != nil {
		return "could not access cache path: " + test
	}
	os.Remove(test)
	return "cache path is accessible at: " + conf.Val(CACHEPATH)
}

func (ConfigChk chk) Amazon(conf Cfg) []string {
	var s []string
	// just checking for amazon config setup so we can alert the user
	for _, v := range awsConfigPath {
		fp, _ := filepath.Abs(filepath.Join(conf.Home(), v))
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			s = append(s, "missing required AWS config file: "+fp)
		} else {
			s = append(s, "found required AWS config file: "+fp)
		}
	}
	return s
}

func (ConfigChk chk) Player(conf Cfg) []string {
	var out []string
	var intern bool
	p := conf.Val(EXECPLAYER)
	if !conf.Exists(EXECPLAYER) {
		out = append(out, "execPlayer not set, will play files internally")
		intern = true
	} else if _, err := exec.LookPath(p); err == nil {
		out = append(out, EXECPLAYER+" found: "+p+", will play files through "+p)
	} else {
		out = append(out, EXECPLAYER+" set but not found: "+p+", will play files internally")
		intern = true
	}
	if intern {
		// buffer check
		if conf.Exists(USERBUFFERSIZE) && conf.Intval(USERBUFFERSIZE) != 0 {
			out = append(out, USERBUFFERSIZE+" set to "+conf.Val(USERBUFFERSIZE))
		}
	}
	return out
}
