package player

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/dustinmj/renotts/coms"
	"github.com/gordonklaus/portaudio"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// multiplier for padding when silence is needed
const paddingMillis = 2000

// intermediate buffer prevents file io calls in audio callback
const interBuffSize = 4000

// minimum intermediate buffer size before we attempt re-fill
const interBuffThresh = 1000

type track struct {
	Decoder  *mpg123.Decoder
	Before   *bytes.Buffer
	After    *bytes.Buffer
	Rate     int64
	Channels int
	Loaded   bool
}

var mpgPlayerQueue = []playerQueueFile{}

//structure for implementing engine interface
type mplayer struct{}

var mpgPlayer = mplayer{}

var done sync.WaitGroup
var mpgplaying bool

func (mpgPlayer mplayer) Play(path string, padB bool, padA bool, player string) error {
	// not interested in player...
	return mpgPlayer.playAudio(path, padB, padA, false)
}

func (mpgPlayer mplayer) playAudio(path string, padB bool, padA bool, fromQueue bool) error {
	if mpgplaying && !fromQueue {
		return errors.New("portaudio is busy")
	}
	mpgplaying = true
	coms.Msg("Playing file: " + path)

	handle, err := mpg123.NewDecoder("TTS")
	if err != nil {
		return err
	}
	defer handle.Close()

	if err = format(handle, path); err != nil {
		mpgplaying = false
		return err
	}

	t := track{
		Decoder: handle}

	// format data
	t.Rate, t.Channels, _ = handle.GetFormat()
	// fill silence
	if padB {
		t.Before = bytes.NewBuffer(getSilence(paddingMillis, t.Rate))
	}
	if padA {
		t.After = bytes.NewBuffer(getSilence(paddingMillis, t.Rate))
	}

	// initialize portaudio only if not from queue
	// when from queue, portaudio is already initialized
	if !fromQueue {
		if err = portaudio.Initialize(); err != nil {
			coms.Msg("Could not initialize portaudio:", err.Error())
			mpgplaying = false
			return err
		}
	}

	// setup channel for capturing kill requests during play
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		portaudio.Terminate() // stop all buffers before exiting
		os.Exit(2)
	}()

	playMPG(&t) // blocking

	// check for queued files
	if len(mpgPlayerQueue) > 0 {
		coms.Msg("Completed Playing. Processing Queue.")
		next := mpgPlayerQueue[0]
		mpgPlayerQueue = mpgPlayerQueue[1:]
		// block for next queued file
		return mpgPlayer.playAudio(next.Path, next.Before, next.After, true)
	}

	// call this directly before closing
	// only terminate after the queue is empty
	if err = portaudio.Terminate(); err != nil {
		coms.Msg(err.Error())
	}
	// sleep momentarily to allow portaudio to finish
	time.Sleep(time.Duration(500) * time.Millisecond)
	// release intterupts
	signal.Reset(os.Interrupt, syscall.SIGTERM)
	// exported Play can now be called
	mpgplaying = false
	coms.Msg("Completed Playing.")

	return nil
}

func (mpgPlayer mplayer) Busy() bool {
	return mpgplaying
}

func (mpgPlayer mplayer) Queue(path string, before bool, after bool) error {
	mpgPlayerQueue = append(mpgPlayerQueue,
		playerQueueFile{
			Path:   path,
			Before: before,
			After:  after})
	return nil
}

// portaudio can hang closing stream,
// we do our best to close out, but
// don't let it block indefinitely
func closeStream(s *portaudio.Stream) {
	s.Abort()
	s.Close()
}

func playMPG(t *track) {
	// load default output device
	out := findOutputDevice()
	if out == nil {
		coms.Msg("Could not determine default output device.")
		unload(t)
		return
	}
	// create parameters
	p := portaudio.HighLatencyParameters(nil, out)
	p.Output.Channels = t.Channels
	p.SampleRate = float64(t.Rate)
	// allow portaudio to decide buffer size
	p.FramesPerBuffer = portaudio.FramesPerBufferUnspecified
	// create stream
	stream, err := portaudio.OpenStream(p, t.playCallback)
	defer closeStream(stream)
	if err != nil {
		coms.Msg(err.Error())
		return
	}

	// ASYNC Stream is the only way to allow PA to determine
	// buffer size. This prevents having to hard-code it which
	// causes issues depending on playback hardware
	done.Add(1)
	t.Loaded = true
	if err = stream.Start(); err != nil {
		coms.Msg(err.Error())
		unload(t)
		return
	}
	done.Wait()
}

func (t *track) playCallback(out []int16) {
	// create output bytes
	audio := make([]byte, len(out)*2)
	if !t.Loaded {
		// since portaudio bindings will not allow
		// us to return values from callback to close
		// stream, a small delay here prevents
		// close stream from hanging indefinitely
		// when close is called called upstream
		time.Sleep(time.Duration(500) * time.Millisecond)
		return
	}
	// read any before data first
	if t.Before != nil && t.Before.Len() > 0 {
		t.Before.Read(audio)
	} else { // before data empty or alread read
		// attempt to read primary
		r, _ := t.Decoder.Read(audio)
		// if primary is empty, read from After
		if r == 0 { // primary is empty
			if t.After != nil && t.After.Len() > 0 {
				t.After.Read(audio)
			} else { // if we've exhausted all reads, unload
				unload(t)
				return
			}
		}
	}
	// read to output
	if rErr := binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out); rErr != nil {
		coms.Msg(rErr.Error())
		unload(t)
		return
	}
	return
}

func unload(t *track) {
	if t.Loaded {
		t.Loaded = false
		done.Done() // unblock
	}
}

func getSilence(mill int, rate int64) []byte {
	return make([]byte, int(float64(rate)*float64(mill/1000)))
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

func findOutputDevice() *portaudio.DeviceInfo {
	// first attempt to find default through pa
	dev, err := portaudio.DefaultOutputDevice()
	if err == nil {
		return dev
	}
	// next, attempt to look for ourselves
	devs, err := portaudio.Devices()
	if err != nil { // if we can't find devices even
		return nil
	}
	// look for nice options
	labels := []string{"default", "pulse", "dmix"}
	if d := searchNames(devs, labels); d != nil {
		return d
	}
	// look for blocking options
	labels = []string{"sysdefault", "spdif", "iec958", "hw"}
	if d := searchNames(devs, labels); d != nil {
		return d
	}
	return nil
}

func searchNames(devs []*portaudio.DeviceInfo, labels []string) *portaudio.DeviceInfo {
	for _, d := range devs {
		for _, l := range labels {
			if strings.ToLower(d.Name) == l && d.MaxOutputChannels > 0 {
				return d
			}
		}
	}
	return nil
}
