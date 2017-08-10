package aud

import (
	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
	"github.com/dustinmj/renotts/com"
	"io"
)

// Play123 execute a play command
func Play123(sF com.Sf) error { // TODO
	com.Msg("Playing file: " + sF.Path)

	mpg123.Initialize()
	defer mpg123.Exit()

	handle, err := mpg123.Open(sF.Path)
	if err != nil {
		return err
	}
	defer handle.Close()

	ao.Initialize()
	defer ao.Shutdown()

	dev := ao.NewLiveDevice(AoSampleFormat(handle))
	defer dev.Close()

	if _, err := io.Copy(dev, handle); err != nil {
		return err
	}
	return nil
}

// AoSampleFormat Get the ao.SampleFormat from the mpg123.Handle
func AoSampleFormat(handle *mpg123.Handle) *ao.SampleFormat {
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
