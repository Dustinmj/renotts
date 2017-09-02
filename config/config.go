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

//Cfg Config interface
type Cfg interface {
	//SetDef - set a default value for config
	SetDef(string, string)
	//SetOverride - override a config value
	SetOverride(string, string)
	// Val retrieve value from config
	Val(string) string
	//Exists does a config key/value Exists
	Exists(string) bool
	//Path gets current config path
	Path() string
	//File gets full path to config file
	File() string
	//User gets current user name
	User() string
	//AppPath gets path to executable
	AppPath() string
	//AppName gets current application name
	AppName() string
	//Home gets current user home directory
	Home() string
	//Cache gets current user home directory
	Cache() string
}

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

type blah struct{}

// conf - main configuration
var conf blah

//AppPath current application path
var appPath string
var appName string

//homeDir maps to user home directory
var homeDir string
var usr *user.User

var configPath string
var defCachePath string

var awsConfigPath = []string{
	".aws/config",
	".aws/credentials"}

func init() {
	// get user
	var err error
	usr, err = user.Current()
	cfg := Get()
	if err != nil {
		coms.Msg("Unable to determine user home directory!!!")
		homeDir = "."
	}
	// get directory info
	homeDir = usr.HomeDir
	appPath, _ = os.Executable()
	appName = os.Args[0]
	// look for amazon config files
	chkAmazonConfig(cfg)
	configPath, _ = filepath.Abs(filepath.Join(homeDir, defRenoFolder))
	defCachePath, _ = filepath.Abs(filepath.Join(configPath, defCacheFolder))
	// setup Viper config path
	setPaths()
	// make sure config file exists
	chkConfigFile(cfg)
	viper.SetConfigName(defConfigName)
	err = viper.ReadInConfig()
	if err == nil {
		coms.Msg("RenoTTS Configuration file loaded:", viper.ConfigFileUsed())
	} else {
		coms.Msg("RenoTTS Configuration file not found, using defaults.")
	}
	// check external player if needed
	chkExecPlayer(cfg)
	// set and check defaults
	setDefs(defPort, defPath, defCachePath)
	chkDefs(cfg)
}

//Get get current config
func Get() Cfg {
	return conf
}

//SetDef - set a default value for config
func (conf blah) SetDef(key string, def string) {
	viper.SetDefault(key, def)
}

//SetOverride - override a config value
func (conf blah) SetOverride(key string, def string) {
	viper.Set(key, def)
}

// Val retrieve value from config
func (conf blah) Val(key string) string {
	return viper.GetString(key)
}

//Exists does a config key/value Exists
func (conf blah) Exists(key string) bool {
	return viper.IsSet(key)
}

//Path gets current config path
func (conf blah) Path() string {
	return configPath
}

//Cache gets current cache path
func (conf blah) Cache() string {
	return conf.Val(CACHEPATH)
}

//Home gets current user home directory
func (conf blah) Home() string {
	return homeDir
}

//File gets full path to config file
func (conf blah) File() string {
	return filepath.Join(configPath, defConfigFile)
}

//User gets current user name
func (conf blah) User() string {
	return usr.Username
}

//AppPath gets path to executable
func (conf blah) AppPath() string {
	return appPath
}

//AppName gets current application name
func (conf blah) AppName() string {
	return appName
}

func setDefs(port string, path string, cachepath string) {
	viper.SetDefault(PORT, port)
	viper.SetDefault(PATH, path)
	viper.SetDefault(CACHEPATH, cachepath)
}

func setPaths() {
	viper.AddConfigPath(configPath)
}

func chkConfigFile(conf Cfg) {
	// check config file to make sure it exists
	cfg := conf.File()
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

func chkAmazonConfig(conf Cfg) {
	// just checking for amazon config setup so we can alert the user
	c := ConfigChk.Amazon(conf)
	for _, t := range c {
		coms.Msg(t)
	}
}

func chkExecPlayer(conf Cfg) {
	// check execplayer settings
	coms.Msg(ConfigChk.Player(conf))
}

// check user input
func chkDefs(conf Cfg) {
	// check cache path to make sure it's writeable
	if _, err := os.Stat(conf.Cache()); os.IsNotExist(err) {
		// attempt to make directory
		if err = os.MkdirAll(conf.Cache(), os.ModePerm); err != nil {
			coms.Msg("Cache directory", conf.Cache(), "not writeable!")
			coms.Exit(73, []byte{})
		}
	}
	// check port to make sure it's correct
	badP := func(conf Cfg) {
		coms.Msg("Invalid port", conf.Val(PORT))
		coms.Exit(78, []byte{})
	}
	re := regexp.MustCompile("[\x3A]?(?P<pnum>\\d{1,5})")
	m := re.FindStringSubmatch(conf.Val(PORT))
	if len(m) < 2 {
		badP(conf)
	}
	mi, err := strconv.Atoi(m[1])
	if err != nil {
		badP(conf)
	}
	conf.SetOverride(PORT, strconv.Itoa(mi))
	// setup path -- we just remove beginning slash
	p := conf.Val(PATH)
	if p[0:1] == "/" {
		p = p[1:]
	}
	conf.SetOverride(PATH, path.Clean(p))
}
