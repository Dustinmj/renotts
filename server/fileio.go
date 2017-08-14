package server

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
	"io"
	"os"
	"path/filepath"
)

//Sf sound file
type Sf struct {
	Q         Aq
	Path      string
	Fname     string
	Pad       padding
	FromCache bool
	ForPlayer bool // future
}

type padding struct {
	Before bool
	After  bool
}

//WriteBuffer - writes buffer stream to file system
//receives an audioQuery com
func WriteBuffer(aQ Aq, rQ *Rq) (Sf, error) {
	coms.Msg("Attempting to copy buffer to file...")
	// create filename
	fN := FileName(rQ)
	// create filepath
	fP, err := FullPath(fN)
	if err != nil {
		return Sf{}, err
	}
	f, err := os.OpenFile(fP, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return Sf{}, err
	}
	// write to file
	_, err = io.Copy(f, *aQ.Buffer)
	if err != nil {
		return Sf{}, err
	}
	return Sf{
		Q:         aQ,
		Path:      fP,
		Fname:     fN,
		FromCache: false}, nil
}

//GetFile - return a Sf for the Rq or error if it doesn't exist yet
func GetFile(rQ *Rq) (Sf, error) {
	fN := FileName(rQ)
	fP, err := FullPath(fN)
	if err != nil {
		return Sf{}, err
	}
	if _, err := os.Stat(fP); os.IsNotExist(err) {
		return Sf{}, err
	}
	return Sf{
		Q:         Aq{},
		Path:      fP,
		Fname:     fN,
		FromCache: true}, nil
}

//FullPath - returns the full path to a file
func FullPath(fN string) (string, error) {
	return filepath.Abs(filepath.Join(config.Val("cachepath"), fN))
}

//FileName - returns the proper filename for caching
func FileName(rQ *Rq) string {
	hasher := md5.New()
	hasher.Write([]byte(rQ.Typ + string(rQ.Body)))
	return hex.EncodeToString(hasher.Sum(nil)) + ".mp3"
}
