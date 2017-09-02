package main

import (
	"fmt"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/file"
	"github.com/dustinmj/renotts/server"
	"os/user"
)

func main() {
	// A simple warning about running RenoTTS as elevated user
	chkElevated()
	// get config
	cfg := config.Get()
	server.Create(cfg.Val(config.PORT), cfg.Val(config.PATH), cfg)
	file.Setup(cfg)
	server.Serve()
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
