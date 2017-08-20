package server

import (
	"bytes"
	"encoding/binary"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/dustinmj/renotts/coms"
	"github.com/gordonklaus/portaudio"
	"os"
	"os/signal"
	"time"
)

const bufferSize = 2024

type dec123 struct {
	Decoder *mpg123.Decoder
}

var mpgPlayer = player{}
var sig chan os.Signal

func (mpgPlayer player) play(sF Sf) error {
	coms.Msg("Playing file: " + sF.Path)

	sig = make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	defer signal.Stop(sig)

	handle, err := mpg123.NewDecoder("TTS")
	if err != nil {
		return err
	}
	defer handle.Close()

	if err = format(handle, sF.Path); err != nil {
		return err
	}

	plays := []*mpg123.Decoder{
		handle}

	// silence
	if isPadded(sF) {
		// shadow lint
		var silence *mpg123.Decoder
		silence, err = mpg123.NewDecoder("SIL")
		if err != nil {
			return err
		}
		defer handle.Close()

		if err = format(silence, silenceFile); err != nil {
			return err
		}
		if sF.Pad.Before {
			plays = append([]*mpg123.Decoder{silence}, plays...)
		}
		if sF.Pad.After {
			plays = append(plays, silence)
		}
	}

	// portaudio can do all types of crazy things if Alsa returns errors
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			coms.Msg("An unrecoverable error ocurred.")
			coms.Exit(74, []byte("Could not play audio file."))
		}
	}()
	// initialize po, capture sdtout since po likes to print
	// all kinds of warnings on some hardware... these don't
	// result in errs but do dump to sdtout, we just ignore
	// them... // TODO is there a better way?
	oldStd := os.Stdout
	os.Stdout = os.Stderr
	defer func() {
		os.Stdout = oldStd
	}()

	if err = portaudio.Initialize(); err != nil {
		coms.Msg("Could not initialize portaudio:", err.Error())
		return nil
	}
	os.Stdout = oldStd
	defer portaudio.Terminate()

	for _, p := range plays {
		if err = playMPG(p); err != nil {
			return err
		}
	}
	return nil
}

func playMPG(p *mpg123.Decoder) error {
	rate, channels, _ := p.GetFormat()
	out := make([]int16, bufferSize)
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	if err != nil {
		return err
	}
	defer stream.Close()
	if err = stream.Start(); err != nil {
		return err
	}
	for {
		audio := make([]byte, 2*len(out))
		read, err := p.Read(audio)
		if read == 0 {
			// we're done... pause for a second
			// to allow buffer to finish Playing
			time.Sleep(time.Duration(1) * time.Second)
			break
		} else if err != nil && err != mpg123.EOF {
			return err
		}
		if err = binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out); err != nil {
			return err
		}
		if err = stream.Write(); err != nil {
			return err
		}
		select {
		case <-sig:
			coms.Exit(74, []byte("Interrupted."))
		default:
		}
	}
	return nil
}

func format(dec *mpg123.Decoder, file string) error {

	if err := dec.Open(file); err != nil {
		return err
	}

	// format info
	rate, channels, _ := dec.GetFormat()
	// don't allow format to vary
	dec.FormatNone()
	dec.Format(rate, channels, mpg123.ENC_SIGNED_16)
	return nil
}

func isPadded(sF Sf) bool {
	return sF.Pad.After || sF.Pad.Before
}
