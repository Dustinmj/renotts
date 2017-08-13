package coms

import (
	"fmt"
	"os"
)

//AppName - name of app
const AppName = "RenoTTS"

//AppVers - version of app
const AppVers = "1.0.0"

// DeviceType - the device type
const DeviceType = "urn:schemas-dustinjorge-com:device:TTSEngine:1"

//Msg - send fmt message with appname:<txt> format
func Msg(txt ...string) {
	str := ""
	for _, t := range txt {
		str += " " + t
	}
	fmt.Println(AppName + ": *" + str)
}

//Exit - send instructions when exiting with error code
func Exit(c int, m []byte) {
	fmt.Println(string(m))
	os.Exit(c)
}
