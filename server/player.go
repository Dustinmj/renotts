package server

import "io"

//Aq audio query
type Aq struct {
	Txt    string
	Typ    string
	Buffer *io.ReadCloser
	Chars  int64
}

type sPlayer interface {
	play(*Sf) error
}

//Serv - structure for implementing engine interface
type player struct{}
