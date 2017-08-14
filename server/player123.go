package server

import (
	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
	"github.com/dustinmj/renotts/coms"
	"io"
)

type play123 struct {
	Handle *mpg123.Handle
	Device *ao.Device
}

var mpgPlayer = player{}

// Play123 execute a play command using mpg123 bindings
func (mpgPlayer player) play(sF Sf) error {
	coms.Msg("Playing file: " + sF.Path)

	mpg123.Initialize()
	defer mpg123.Exit()

	handle, err := mpg123.Open(sF.Path)
	if err != nil {
		return err
	}
	defer handle.Close()
	return mpgPlay(handle, sF)
}

func mpgPlay(handle *mpg123.Handle, sF Sf) error {
	ao.Initialize()
	defer ao.Shutdown()

	dev := ao.NewLiveDevice(aoSampleFormat(handle))
	defer dev.Close()

	plays := []play123{
		play123{
			Device: dev,
			Handle: handle}}
	// if we pad with silence, open a device for that
	if isPadded(sF) {
		sHandle, err := mpg123.Open(silenceFile)
		if err != nil {
			return err
		}
		sdev := ao.NewLiveDevice(aoSampleFormat(sHandle))
		silence := play123{
			Device: sdev,
			Handle: sHandle}
		if sF.Pad.Before {
			plays = append([]play123{silence}, plays...)
		}
		if sF.Pad.After {
			plays = append(plays, silence)
		}
	}
	for _, p := range plays {
		if _, err := io.Copy(p.Device, p.Handle); err != nil {
			return err
		}
	}
	return nil
}

func isPadded(sF Sf) bool {
	return sF.Pad.After || sF.Pad.Before
}

// AoSampleFormat Get the ao.SampleFormat from the mpg123.Handle
func aoSampleFormat(handle *mpg123.Handle) *ao.SampleFormat {
	const bitsPerByte = 8

	rate, channels, encoding := handle.Format()

	return &ao.SampleFormat{
		BitsPerSample: handle.EncodingSize(encoding) * bitsPerByte,
		Rate:          int(rate),
		Channels:      channels,
		ByteFormat:    ao.FormatNative,
		Matrix:        nil,
	}
}
