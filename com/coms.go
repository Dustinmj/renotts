package com

import "io"

/*Eng tts engine interface
Query returns the sound file and http response object
SetDefs allows engine to write it's own defaults, these can then be used and will
be honored if they exist in the main config file from user
Caches return true if the Sf file returned from query will need to be played,
it will then be passed to the audio player (mp3);
return false if the engine handles playing the file internally
*/
type Eng interface {
	Query(*Rq) (Sf, Rsp)
	SetDefs()
	Caches() bool
}

//Param http request structure
type Param struct {
	Text, Voice, SampleRate string
}

//Rq http request
type Rq struct {
	Typ   string
	Param Param
	Body  []byte
}

//Rsp http response
type Rsp struct {
	Msg   string
	Err   error
	Code  int
	Heads map[string]string
}

//Aq audio query
type Aq struct {
	Txt    string
	Typ    string
	Buffer *io.ReadCloser
	Chars  int64
}

//Sf sound file
type Sf struct {
	Q         Aq
	Path      string
	Fname     string
	FromCache bool
	ForPlayer bool
}
