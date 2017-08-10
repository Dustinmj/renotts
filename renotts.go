package main

import (
	"github.com/dustinmj/renotts/serv"
	"github.com/dustinmj/renotts/upnp"
)

func main() {
	upnp.Create()
	serv.Create()
}
