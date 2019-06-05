package opsman

import (
	"io"
	"io/ioutil"
)

type liveDiscarder struct {
}

func (ld liveDiscarder) Write(b []byte) (int, error) {
	return ioutil.Discard.Write(b)
}

func (ld liveDiscarder) Start() {
}

func (ld liveDiscarder) Stop() {

}
func (ld liveDiscarder) Flush() error {
	return nil
}

type progressBarDiscarder struct {
}

func (pbd progressBarDiscarder) NewProxyReader(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}

func (pbd progressBarDiscarder) Start() {
}

func (pbd progressBarDiscarder) Finish() {
}

func (pbd progressBarDiscarder) SetTotal64(_ int64) {
}

func (pbd progressBarDiscarder) Reset() {
}

func (pbd progressBarDiscarder) SetOutput(_ io.Writer) {
}
