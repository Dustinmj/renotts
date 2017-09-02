package main

import (
	"fmt"
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/file"
	"github.com/dustinmj/renotts/server"
	"os"
	"os/user"
)

func main() {
	// A simple warning about running RenoTTS as elevated user
	chkElevated()
	// get config
	cfg := config.Get()
	server.Create(cfg)
	file.Setup(cfg)
	err := server.Serve() // blocking
	if err != nil {
		coms.Msg(err.Error())
		os.Exit(2)
	}
}

func chkElevated() {
	warn := "!!!!! RenoTTS appears to be running as root user. This is not recommended. !!!!!"
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
