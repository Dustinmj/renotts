package serv

import (
	"errors"
	"github.com/dustinmj/renotts/com"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/service"
)

func enGet(t string) (e com.Eng, er error) {
	for n, s := range service.AvailServs {
		if n == t {
			s.SetDefs()
			return s, nil
		}
	}
	return nil, errors.New(config.Err["NoService"])
}
