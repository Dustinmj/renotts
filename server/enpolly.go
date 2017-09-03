package server

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/file"
)

// error messages
const (
	errAWSFailed = "aws request failed"
)

//Polly service namespace
var Polly = engine{}

// default configurations for polly
func (Polly engine) SetDefs() {}
func (Polly engine) Caches() bool {
	return true
}

func (Polly engine) Query(req *request, cfg config.Cfg) (*string, error) {
	sF, err := file.GetFile(req.Unique, cfg.Cache()) // checks for cached file
	if err != nil {
		sF, err = awsRequest(req, cfg.Val(config.AWSPROFILE), cfg.Cache())
		if err != nil {
			return nil, err
		}
	}
	return sF, nil
}

func awsRequest(req *request, prf string, cache string) (*string, error) {
	coms.Msg("Sending aws request...")
	format := "mp3"
	sample := req.Param.SampleRate
	voice := req.Param.Voice
	text := req.Param.Text
	var sess *session.Session
	if len(prf) > 0 {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable, Profile: prf,
		}))
	} else {
		// Force enable Shared Config support
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}
	client := polly.New(sess)
	params := polly.SynthesizeSpeechInput{
		OutputFormat: &format,
		SampleRate:   &sample,
		VoiceId:      &voice,
		Text:         &text}
	to, from := client.SynthesizeSpeechRequest(&params)
	err := to.Send()
	if err != nil {
		coms.Msg("AWS Request Failed.")
		return nil, err
	}
	defer from.AudioStream.Close()
	sF, err := file.WriteBuffer(&from.AudioStream, req.Unique, cache)
	if err != nil {
		coms.Msg("Error writing file, check cache path settings.")
		return nil, err
	}
	return sF, nil
}
