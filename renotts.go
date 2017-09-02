package main

import (
	"fmt"
	"github.com/dustinmj/renotts/server"
	"os/user"
)

func main() {
	// A simple warning about running RenoTTS as elevated user
	chkElevated()
	server.Create()
}

func chkElevated() {
	warn := "!!!!! RenoTTS appears to be running as root user. This is not recommended. RenoTTS should not be run as root user. !!!!!"
	user, err := user.Current()
	if err != nil {
		return
	}
	if user.Uid == "0" || user.Username == "root" {
		fmt.Println("")
		fmt.Println(warn)
		fmt.Println("")
	}
}
