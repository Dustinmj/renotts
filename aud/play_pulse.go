package aud

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/dustinmj/renotts/com"
	"github.com/gordonklaus/portaudio"
	"io"
	"os"
	"os/signal"
)

//PlayPulse execute a play command throug pulseaudio bindings
func PlayPulse(sF com.Sf) error {
	if len(sF.Path) < 2 {
		return errors.New("missing required argument:  input file name")
	}
	com.Msg("Playing " + sF.Path)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	defer signal.Reset()

	f, err := os.Open(sF.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, data, err := readChunk(f)
	if err != nil {
		return err
	} /*
		fmt.Println(id.String())
		if id.String() != "FORM" {
			return errors.New("bad file format")
		}
		_, err = data.Read(id[:])
		if err != nil {
			return err
		}
		if id.String() != "AIFF" {
			return errors.New("bad file format")
		}*/
	var c commonChunk
	var audio io.Reader
	for {
		id, chunk, errs := readChunk(data)
		if errs == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch id.String() {
		case "COMM":
			if errs := binary.Read(chunk, binary.BigEndian, &c); err != nil {
				return errs
			}
		case "SSND":
			chunk.Seek(8, 1) //ignore offset and block
			audio = chunk
		default:
			fmt.Printf("ignoring unknown chunk '%s'\n", id)
		}
	}

	//aws 22050 sample rate, mono, 32 bit

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, 22050, len(out), &out)
	if err != nil {
		return err
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return err
	}
	defer stream.Stop()
	for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err := stream.Write(); err != nil {
			return err
		}
		select {
		case <-sig:
			return nil
		default:
		}
	}
	return nil
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

//ID - idk
type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}
