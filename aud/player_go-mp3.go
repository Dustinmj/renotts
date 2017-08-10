package aud

import (
	"bufio"
	"fmt"
	"github.com/dustinmj/renotts/com"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"os"
)

/*func init() {
	sF := com.Sf{
		Path: "/home/dustin/easyttscache/classic.mp3"}
	Play(sF)
}*/ // testing
type readCloser struct {
	io.Reader
	io.Closer
}

//PlayHaj - plays a Sf Sound File using hajimehoshi's go-mp3,
//note - does not support mpeg v2 so pretty worthless to this app
func PlayHaj(sF com.Sf) error {
	fmt.Println("Aud: Playing", sF.Path)

	f, err := os.Open(sF.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewReader(f)
	d, err := mp3.NewDecoder(&readCloser{b, f})
	if err != nil {
		return err
	}
	defer d.Close()

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
