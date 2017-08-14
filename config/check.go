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
	var s []string
	s = append(s, ConfigChk.Config()...)
	s = append(s, ConfigChk.Cache())
	s = append(s, ConfigChk.Player())
	s = append(s, ConfigChk.Amazon()...)
	return s
}

func (ConfigChk chk) Config() []string {
	var s []string
	// check config file to make sure it exists
	cfg := fullConfigPath()
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		s = append(s, "missing config file: "+cfg)
	}
	s = append(s, "found renotts config file: "+cfg)
	// check config values
	vals := viper.AllKeys()
	for _, k := range vals {
		s = append(s, "** current Configuration -> "+k+": "+Val(k))
	}
	return s
}

func (ConfigChk chk) Cache() string {
	// check config file to make sure it exists
	test := filepath.Join(Val(CACHEPATH), "test.txt")
	if _, err := os.OpenFile(test, os.O_RDWR|os.O_CREATE, 0755); err != nil {
		return "could not access cache path: " + test
	}
	os.Remove(test)
	return "cache path is accessible at: " + Val(CACHEPATH)
}

func (ConfigChk chk) Amazon() []string {
	var s []string
	// just checking for amazon config setup so we can alert the user
	for _, v := range awsConfigPath {
		fp, _ := filepath.Abs(filepath.Join(homeDir, v))
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			s = append(s, "missing AWS config file: "+fp)
		} else {
			s = append(s, "found required AWS config file: "+fp)
		}
	}
	return s
}

func (ConfigChk chk) Player() string {
	if !viper.IsSet(EXECPLAYER) {
		return "execPlayer not set, will play files internally"
	}
	// if execplayer option is set, make sure we have access to it
	p := Val(EXECPLAYER)
	if _, err := exec.LookPath(p); err == nil {
		return "execPlayer found: " + p + ", will play files through " + p
	}
	return "execPlayer set but not found: " + p + ", will play files internally"
}
