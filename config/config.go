package config

import (
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/tmplt"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

// Config file strings
const (
	//AWSPROFILE blah
	AWSPROFILE = "aws-config-profile"
	//EXECPLAYER blah
	EXECPLAYER = "exec-player"
	//CACHEPATH blah
	CACHEPATH = "cache-path"
	//PORT blah
	PORT = "port"
	//PATH blah
	PATH = "path"
)

// config defaults
const (
	defCacheFolder = "cache"
	defRenoFolder  = ".renotts"
	defConfigName  = "renotts"
	defPort        = "0"
	defPath        = "tts"
	defConfigFile  = "renotts.toml"
	defAWSProfile  = "default"
	defExecPlayer  = "mplayer"
)

//AppPath current application path
var AppPath string

//homeDir maps to user home directory
var homeDir string
var usr *user.User

var configPath string
var defCachePath string

//SilenceFile silence file full path
var SilenceFile string

//Smp3 Silence mp3 file
const Smp3 = "Silence.mp3"

var awsConfigPath = []string{
	".aws/config",
	".aws/credentials"}

func init() {
	// get user
	var err error
	usr, err = user.Current()
	if err != nil {
		coms.Msg("Unable to determine user home directory!!!")
		homeDir = "."
	}
	homeDir = usr.HomeDir
	AppPath, _ = os.Executable()
	chkAmazonConfig()
	configPath, _ = filepath.Abs(filepath.Join(homeDir, defRenoFolder))
	defCachePath, _ = filepath.Abs(filepath.Join(configPath, defCacheFolder))
	setPaths()
	chkConfigFile()
	viper.SetConfigName(defConfigName)
	err = viper.ReadInConfig()
	if err == nil {
		coms.Msg("RenoTTS Configuration file loaded:", viper.ConfigFileUsed())
	} else {
		coms.Msg("RenoTTS Configuration file not found, using defaults.")
	}
	chkExecPlayer()
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

//Exists does a config key/value Exists
func Exists(key string) bool {
	return viper.IsSet(key)
}

//Path gets current config path
func Path() string {
	return configPath
}

//File gets full path to config file
func File() string {
	return filepath.Join(configPath, defConfigFile)
}

//User gets current user name
func User() string {
	return usr.Username
}

func setDefs() {
	viper.SetDefault(PORT, defPort)
	viper.SetDefault(PATH, defPath)
	viper.SetDefault(CACHEPATH, defCachePath)
}

func setPaths() {
	viper.AddConfigPath(configPath)
}

func chkConfigFile() {
	// check config file to make sure it exists
	cfg := fullConfigPath()
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		// attempt to make directory
		if err = os.MkdirAll(configPath, os.ModePerm); err != nil {
			coms.Msg("Could not create config directory:", configPath)
			coms.Exit(73, []byte{})
		}
		// attempt to create base config file
		err := createConfig(cfg)
		if err != nil {
			coms.Msg("Could not create config skeleton:", cfg)
		}
	}
}

func createConfig(path string) error {
	// attempt to create base config file
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	dat := tmplt.ConfigData{
		Awsconfigprofile: defOption(AWSPROFILE, defAWSProfile, true),
		Port:             defOption(PORT, defPort, false),
		Path:             defOption(PATH, defPath, false),
		Execplayer:       defOption(EXECPLAYER, defExecPlayer, true),
		Cachepath:        defOption(CACHEPATH, defCachePath, false)}
	tmplt.ParseF(f, tmplt.ConfigFl, dat)
	return nil
}

func defOption(label string, value string, commented bool) string {
	s := ""
	if commented {
		s += "#"
	}
	s += label + "=\"" + value + "\""
	return s
}

func fullConfigPath() string {
	return filepath.Join(configPath, defConfigFile)
}

func chkAmazonConfig() {
	// just checking for amazon config setup so we can alert the user
	c := ConfigChk.Amazon()
	for _, t := range c {
		coms.Msg(t)
	}
}

func chkExecPlayer() {
	// check execplayer settings
	coms.Msg(ConfigChk.Player())
}

// check user input
func chkDefs() {
	// check cache path to make sure it's writeable
	if _, err := os.Stat(Val(CACHEPATH)); os.IsNotExist(err) {
		// attempt to make directory
		if err = os.MkdirAll(Val(CACHEPATH), os.ModePerm); err != nil {
			coms.Msg("Cache directory", Val(CACHEPATH), "not writeable!")
			coms.Exit(73, []byte{})
		}
	}
	// check port to make sure it's correct
	badP := func() {
		coms.Msg("Invalid port", Val(PORT))
		coms.Exit(78, []byte{})
	}
	re := regexp.MustCompile("[\x3A]?(?P<pnum>\\d{1,5})")
	m := re.FindStringSubmatch(Val(PORT))
	if len(m) < 2 {
		badP()
	}
	mi, err := strconv.Atoi(m[1])
	if err != nil {
		badP()
	}
	SetOverride(PORT, strconv.Itoa(mi))
	// setup path -- we just remove beginning slash
	p := Val(PATH)
	if p[0:1] == "/" {
		p = p[1:]
	}
	SetOverride(PATH, path.Clean(p))
}
