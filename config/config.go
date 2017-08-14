package config

import (
	"github.com/dustinmj/renotts/coms"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

var awsConfigPath = []string{
	".aws/config",
	".aws/credentials"}

//HomeDir - maps to user home directory
var HomeDir string

var defConfigPath string
var defCachePath string

func init() {
	HomeDir, _ = homedir.Dir()
	chkAmazonConfig()
	defConfigPath, _ = filepath.Abs(filepath.Join(HomeDir, ".renotts"))
	defCachePath, _ = filepath.Abs(filepath.Join(defConfigPath, "cache"))
	setPaths()
	chkConfigFile()
	viper.SetConfigName("renotts")
	err := viper.ReadInConfig()
	if err == nil {
		coms.Msg("Configuration file loaded:", viper.ConfigFileUsed())
	} else {
		coms.Msg("Configuration file not found, using defaults.")
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
	viper.SetDefault("cachepath", defCachePath)
}

func setPaths() {
	viper.AddConfigPath(".")
	viper.AddConfigPath(defConfigPath)
}

func chkConfigFile() {
	// check config file to make sure it exists
	cfg := filepath.Join(defConfigPath, "renotts.toml")
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		// attempt to make directory
		if err = os.MkdirAll(defConfigPath, os.ModePerm); err != nil {
			coms.Msg("Could not create config directory:", defConfigPath)
			coms.Exit(73, []byte{})
		}
		err := ioutil.WriteFile(cfg, defConfig, 0744)
		if err != nil {
			coms.Msg("Could not create config skeleton:", cfg)
		}
	}
}

func chkAmazonConfig() {
	// just checking for amazon config setup so we can alert the user
	for _, v := range awsConfigPath {
		fp, _ := filepath.Abs(filepath.Join(HomeDir, v))
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			coms.Msg("Note: AWS configuration missing:", fp)
		}
	}
}

// check user input
func chkDefs() {
	// check cache path to make sure it's writeable
	if _, err := os.Stat(Val("cachepath")); os.IsNotExist(err) {
		// attempt to make directory
		if err = os.MkdirAll(Val("cachepath"), os.ModePerm); err != nil {
			coms.Msg("Cache directory", Val("cachepath"), "not writeable!")
			coms.Exit(73, []byte{})
		}
	}
	// check port to make sure it's correct
	badP := func() {
		coms.Msg("Invalid port", Val("port"))
		coms.Exit(78, []byte{})
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
