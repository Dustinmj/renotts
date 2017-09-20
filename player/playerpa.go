package player

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/dustinmj/renotts/coms"
	"github.com/gordonklaus/portaudio"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// max expected sample rate, this is used
// to approximate wait time during blocking
// playback
const maxSampleRate = 22050

// number of milliseconds to pad for silence
const paddingMillis = 1000

// default size for intermediate read buffer
const rdBuffSize = 10000

// maximum buffer size user can set
const maxBufferSize = 60000

// min buffer size user can set
const minBufferSize = 10

// user configured device type name
// allows user to create specific device
// type for renotts (e.g. in ~/.asoundrc),
// this will be chosen over any other, if it exists.
const userDeviceTypeName = "renotts"

type track struct {
	Decoder  *mpgBuff
	Before   *bytes.Buffer
	After    *bytes.Buffer
	Rate     int64
	Channels int
	Loaded   bool
}

var mpgPlayerQueue = []playerQueueFile{}

//structure for implementing engine interface
type mplayer struct {
	userBufferSize int
}

var done sync.WaitGroup
var mpgplaying bool

//player string
func (mpgPlayer *mplayer) Play(path string, padB bool, padA bool) error {
	return mpgPlayer.playAudio(path, padB, padA, false)
}

func (mpgPlayer *mplayer) playAudio(path string, padB bool, padA bool, fromQueue bool) error {
	if mpgplaying && !fromQueue {
		return errors.New("portaudio is busy")
	}
	mpgplaying = true
	coms.Msg("Playing file: " + path)

	async := (mpgPlayer.BSize() == 0)

	handle, err := mpg123.NewDecoder("TTS")
	if err != nil {
		return err
	}
	defer handle.Close()

	if err = format(handle, path); err != nil {
		mpgplaying = false
		return err
	}
	var t track
	if !async {
		t = track{}
	} else {
		t = track{
			Decoder: &mpgBuff{
				fileDec: handle,
				size:    rdBuffSize}}
		// prepare the intermediate buffer
		t.Decoder.Prepare()
	}
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
		if err = paInit(); err != nil {
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
		close(sig)
		portaudio.Terminate()
		os.Exit(2)
	}()

	if async {
		// blocks until complete and stream closed or err
		mpgPlayer.playMPG(&t)
	} else {
		// blocks until complete and stream closed or err
		mpgPlayer.syncPlayMPG(&t, handle)
	}

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
	// release intterupts
	signal.Reset(os.Interrupt, syscall.SIGTERM)
	// exported Play can now be called
	mpgplaying = false
	coms.Msg("Completed Playing.")

	return nil
}

func paInit() error {
	return portaudio.Initialize()
}

func (mpgPlayer *mplayer) Busy() bool {
	return mpgplaying
}

// get proper buffer size
func (mpgPlayer *mplayer) BSize() int {
	// check buffer size configuration is within reasonable limits
	bs := mpgPlayer.userBufferSize
	if bs > maxBufferSize {
		bs = maxBufferSize
	} else if bs < minBufferSize && bs > 0 {
		bs = minBufferSize
	}
	if bs > 0 {
		// ensure divisible by 2
		bs += bs % 2
	}
	return bs
}

func (mpgPlayer *mplayer) Queue(path string, before bool, after bool) error {
	mpgPlayerQueue = append(mpgPlayerQueue,
		playerQueueFile{
			Path:   path,
			Before: before,
			After:  after})
	return nil
}

func closeStream(s *portaudio.Stream) {
	c := make(chan struct{}, 1)
	go func() {
		s.Close()
		close(c)
	}()
	<-c
}

func (mpgPlayer *mplayer) playMPG(t *track) {
	p, err := t.paParams(mpgPlayer.BSize())
	if err != nil {
		coms.Msg(err.Error())
		return
	}
	// create stream
	stream, err := portaudio.OpenStream(p, t.playCallback)
	if err != nil {
		coms.Msg(err.Error())
		return
	}
	defer closeStream(stream)

	// ASYNC Stream is the only way to allow PA to determine
	// buffer size. This prevents having to hard-code it which
	// causes issues depending on playback hardware.
	// if user specifies buffer size in config, however, we
	// used that in synchronous playback.
	done.Add(1)
	t.Loaded = true
	if err = stream.Start(); err != nil {
		coms.Msg(err.Error())
		t.unload(0)
		return
	}
	done.Wait()
	stream.Stop()
}

func (t *track) paParams(frames int) (portaudio.StreamParameters, error) {
	// load default output device
	out := findOutputDevice()
	if out == nil {
		return portaudio.StreamParameters{}, errors.New("unable to determine output device")
	}
	// create parameters
	p := portaudio.HighLatencyParameters(nil, out)
	p.Output.Channels = t.Channels
	if frames > 0 {
		// user defined frames per buffer...
		p.FramesPerBuffer = frames
	} else {
		// allow portaudio to decide buffer size, don't 'have' to set this, but... idiomatic
		p.FramesPerBuffer = portaudio.FramesPerBufferUnspecified
	}
	p.SampleRate = float64(t.Rate)
	return p, nil
}

func (t *track) playCallback(out []int16) {
	if !t.Loaded {
		return
	}
	// create output bytes
	audio := make([]byte, len(out)*2)
	r, done := t.readAudio(audio, t.Decoder)
	//
	if done {
		go t.unload(r)
		return
	}
	// read to output
	if rErr := binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out); rErr != nil {
		coms.Msg(rErr.Error())
		go t.unload(0)
		return
	}
	return
}

// block for remaining seconds based on bytes read
// from buffer
func (t *track) wait(r int) {
	<-time.After(time.Second * time.Duration(r/(2*int(t.Rate))))
}

func (t *track) unload(r int) {
	if t.Loaded {
		t.Loaded = false
		// wait for stream to complete
		t.wait(r)
		done.Done() // unblock
	}
}

func (t *track) readAudio(audio []byte, rdr io.Reader) (int, bool) {
	var r int
	done := false
	// read any before data first
	if t.Before != nil && t.Before.Len() > 0 {
		r, _ = t.Before.Read(audio)
	} else {
		// attempt to read primary
		r, _ = rdr.Read(audio)
		// check who determines EOF
		if t.After != nil && t.After.Len() > 0 {
			// if primary is empty and after is not, after is EOF mark
			if r == 0 {
				r, _ = t.After.Read(audio)
				if r < len(audio) {
					done = true
				}
			}
		} else if r < len(audio) {
			done = true
		}
	}
	// if after is not defined, we use primary as EOF mark
	return r, done
}

func (mpgPlayer *mplayer) syncPlayMPG(t *track, h *mpg123.Decoder) {
	p, err := t.paParams(mpgPlayer.BSize())
	if err != nil {
		coms.Msg(err.Error())
		return
	}
	// create output buffer
	out := make([]int16, mpgPlayer.BSize())
	// create stream
	stream, err := portaudio.OpenStream(p, &out)
	if err != nil {
		coms.Msg(err.Error())
		return
	}
	stream.Start()
	defer stream.Stop()
	defer closeStream(stream)
	cnt := 0
	var tm time.Time
	for {
		audio := make([]byte, 2*len(out))
		r, done := t.readAudio(audio, h)
		cnt += r
		err := binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, &out)
		if err != nil {
			coms.Msg(err.Error())
			return
		}
		if err = stream.Write(); err != nil {
			coms.Msg(err.Error())
			return
		}
		// now we've written bytes to the hardware,
		// start expectation clock
		if tm.IsZero() {
			tm = time.Now()
		}
		if done {
			// since we're stuffing the buffer, give stream a chance to complete
			expected := time.Duration(cnt/(2*int(t.Rate))) * time.Second
			actual := time.Since(tm)
			<-time.After(expected - actual)
			// due to hardware buffering, time may vary...
			// wait for 1 more full buffer at max samplerate just to be sure
			// then add "just a little bit" (very accurate)
			rt := float32(t.Rate)/float32(maxSampleRate) + .125
			t.wait(int(float32(len(audio)) * rt))
			break
		}
	}
}

func getSilence(mill int, rate int64) []byte {
	r := float32(int(rate)*2) * float32(mill/1000)
	return make([]byte, int(r))
}

func format(dec *mpg123.Decoder, file string) error {
	if err := dec.Open(file); err != nil {
		return err
	}
	// format info
	rate, chans, _ := dec.GetFormat()
	// don't allow format to vary
	dec.FormatNone()
	dec.Format(rate, chans, mpg123.ENC_SIGNED_16)
	return nil
}

func findOutputDevice() *portaudio.DeviceInfo {
	// first, prefer user defined options
	devs, _ := portaudio.Devices()
	if devs != nil {
		if d := searchNames(devs, []string{userDeviceTypeName}); d != nil {
			return d
		}
	}
	// second, attempt to find default through pa
	dev, err := portaudio.DefaultOutputDevice()
	if err == nil {
		return dev
	}
	// next, attempt to look for ourselves
	if devs != nil {
		// look for nice options
		labels := []string{"default", "pulse", "dmix"}
		if d := searchNames(devs, labels); d != nil {
			return d
		}
		// look for other options
		labels = []string{"sysdefault", "spdif", "iec958", "hw"}
		if d := searchNames(devs, labels); d != nil {
			return d
		}
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
