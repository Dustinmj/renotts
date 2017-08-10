package config

import (
	"github.com/dustinmj/renotts/com"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"path"
	"regexp"
	"strconv"
)

//HomeDir - maps to user home directory
var HomeDir string

func init() {
	HomeDir, _ = homedir.Dir()
	setPaths()
	viper.SetConfigName("renotts")
	err := viper.ReadInConfig()
	if err == nil {
		com.Msg("Configuration file loaded: " + viper.ConfigFileUsed())
	} else {
		com.Msg("Configuration file not found, using defaults.")
	}
	// check port
	setDefs()
	chkDefs()
}

//SetDef - set a default value for config
func SetDef(key string, def string) {
	viper.SetDefault(key, def)
}

//SetOverride - override a config value
func SetOverride(key string, def string) {
	viper.Set(key, def)
}

// Val retrieve value from config
func Val(key string) string {
	return viper.GetString(key)
}

func setDefs() {
	viper.SetDefault("port", "0")
	viper.SetDefault("path", "tts")
	viper.SetDefault("cachepath", path.Clean("."))
}

func setPaths() {
	viper.AddConfigPath(".")
	viper.AddConfigPath(HomeDir + "/.renotts/")
}

// check user input
func chkDefs() {
	// check cache path to make sure it's writeable
	if _, err := os.Stat(Val("cachepath")); os.IsNotExist(err) {
		// attempt to make directory
		if err = os.MkdirAll(Val("cachepath"), os.ModePerm); err != nil {
			com.Msg("Cache directory", Val("cachepath"), "not writeable!")
			com.Exit(73, []byte{})
		}
	}
	// check port to make sure it's correct
	badP := func() {
		com.Msg("!!! Invalid port", Val("port"))
		com.Exit(78, []byte{})
	}
	re := regexp.MustCompile("[\x3A]?(?P<pnum>\\d{1,5})")
	m := re.FindStringSubmatch(Val("port"))
	if len(m) < 2 {
		badP()
	}
	mi, err := strconv.Atoi(m[1])
	if err != nil {
		badP()
	}
	SetOverride("port", strconv.Itoa(mi))
	// setup path -- we just remove beginning slash
	p := Val("path")
	if p[0:1] == "/" {
		p = p[1:]
	}
	SetOverride("path", path.Clean(p))
}
