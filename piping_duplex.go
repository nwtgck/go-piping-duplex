package piping_duplex

import (
	"fmt"
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
		// TODO: use url join
		res, err := http.Get(fmt.Sprintf("%s/%s", server, peerPath))
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
		// TODO: use url join
		_, err := http.Post(fmt.Sprintf("%s/%s", server, selfPath), contentType, input)
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
