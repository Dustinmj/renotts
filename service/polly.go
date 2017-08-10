package service

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/dustinmj/renotts/com"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/fs"
	"net/http"
)

//Serv - structure for implementing server interface
type serv struct{}

//Polly service namespace
var Polly = serv{}

// default configurations for polly
func (Polly serv) SetDefs() {}
func (Polly serv) Caches() bool {
	return true
}

func (Polly serv) Query(req *com.Rq) (com.Sf, com.Rsp) {
	rC := http.StatusOK // success, from cache

	sF, err := fs.GetFile(req) // checks for cached file
	if err != nil {
		sF, err = awsRequest(req)
		rC = http.StatusCreated // reset content
	}

	heads := map[string]string{"Via": polly.ServiceName}

	var msg string
	var rsCd int
	if err == nil {
		rsCd = rC
		msg = "Query Successful"
	} else {
		rsCd = http.StatusInternalServerError
		msg = "AWS Polly query failed, see log for more details."
	}

	return sF, com.Rsp{
		Msg:   msg,
		Err:   err,
		Code:  rsCd,
		Heads: heads}
}

func awsRequest(rQ *com.Rq) (com.Sf, error) {
	com.Msg("Sending aws request...")
	format := "mp3"
	sample := rQ.Param.SampleRate
	voice := rQ.Param.Voice
	text := rQ.Param.Text
	// look for aws profile settings
	prf := config.Val("aws-config-profile")
	var sess *session.Session
	if len(prf) > 0 {
		// Force enable Shared Config support
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
		com.Msg("AWS Request Failed: " + err.Error())
		return com.Sf{}, err
	}
	defer from.AudioStream.Close()
	aQ := com.Aq{
		Txt:    text,
		Typ:    rQ.Typ,
		Chars:  *from.RequestCharacters,
		Buffer: &from.AudioStream,
	}
	sF, err := fs.WriteBuffer(aQ, rQ)
	if err != nil {
		com.Msg("Error writing file, check cache path settings.")
		return com.Sf{}, err
	}
	return sF, nil
}
