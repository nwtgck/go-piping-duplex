package piping_duplex

import (
	"github.com/nwtgck/go-piping-duplex/util"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Wait(server string, selfId string, peerId string) error {
	return Duplex(server, selfId, peerId, strings.NewReader("OK"), ioutil.Discard)
}

func Duplex(server string, selfPath string, peerPath string, input io.Reader, output io.Writer) error {
	c := make(chan error)
	go func() {
		url, err := util.UrlJoin(server, peerPath)
		if err != nil {
			c <- err
			return
		}
		res, err := http.Get(url)
		if err != nil {
			c <- err
			return
		}
		_, err = io.Copy(output, res.Body)
		c <- err
	}()
	go func() {
		// TODO: hard code
		contentType := "application/octet-stream"
		url, err := util.UrlJoin(server, selfPath)
		if err != nil {
			c <- err
			return
		}
		_, err = http.Post(url, contentType, input)
		c <- err
	}()
	var err error
	err = <- c
	if err != nil {
		return err
	}
	err = <- c
	if err != nil {
		return err
	}
	return nil
}
