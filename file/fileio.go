package file

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/dustinmj/renotts/coms"
	"io"
	"os"
	"path/filepath"
)

//WriteBuffer - writes buffer stream to file system
func WriteBuffer(buffer *io.ReadCloser, unique []byte, cache string) (*string, error) {
	coms.Msg("Attempting to copy buffer to file...")
	// create filename
	fN := name(unique)
	// create filepath
	fP, err := FullPath(fN, cache)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(fP, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	// write to file
	_, err = io.Copy(f, *buffer)
	if err != nil {
		return nil, err
	}
	return &fP, nil
}

//GetFile - return a Sf for the Rq or error if it doesn't exist yet
func GetFile(unique []byte, cache string) (*string, error) {
	fN := name(unique)
	fP, err := FullPath(fN, cache)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fP); os.IsNotExist(err) {
		return nil, err
	}
	return &fP, nil
}

//FullPath - returns the full path to a file
func FullPath(fN string, cache string) (string, error) {
	return filepath.Abs(filepath.Join(cache, fN))
}

//Name - returns the proper filename for caching
func name(unique []byte) string {
	hash := md5.Sum(unique)
	return hex.EncodeToString(hash[:]) + ".mp3"
}
