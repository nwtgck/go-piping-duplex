package piping_duplex

import (
	"github.com/nwtgck/go-piping-duplex/util"
	"io"
	"net/http"
)

func Duplex(server string, selfPath string, peerPath string, r io.Reader) (io.Reader, error) {
	postUrl, err := util.UrlJoin(server, selfPath)
	if err != nil {
		return nil, err
	}
	// TODO: hard code
	contentType := "application/octet-stream"
	_, err = http.Post(postUrl, contentType, r)
	if err != nil {
		return nil, err
	}
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
