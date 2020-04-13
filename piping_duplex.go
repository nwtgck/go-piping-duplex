package piping_duplex

import (
	"github.com/nwtgck/go-piping-duplex/util"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Wait(server string, selfId string, peerId string) error {
	r, err :=  Duplex(server, selfId, peerId, strings.NewReader("OK"))
	if err != nil {
		return err
	}
	_, err = io.Copy(ioutil.Discard, r)
	return err
}


func Duplex(server string, selfPath string, peerPath string, r io.Reader) (io.Reader, error) {
	postUrl, err := util.UrlJoin(server, selfPath)
	if err != nil {
		return nil, err
	}
	go func() {
		// TODO: hard code
		contentType := "application/octet-stream"
		_, err = http.Post(postUrl, contentType, r)
		if err != nil {
			panic(err)
		}
	}()
	getUrl, err := util.UrlJoin(server, peerPath)
	if err != nil {
		return nil, err
	}
	res, err := http.Get(getUrl)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
