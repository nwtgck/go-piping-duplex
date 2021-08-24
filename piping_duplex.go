package piping_duplex

import (
	"github.com/nwtgck/go-piping-duplex/util"
	"github.com/pkg/errors"
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
	postRes, err := http.Post(postUrl, contentType, r)
	if err != nil {
		return nil, err
	}
	if postRes.StatusCode != 200 {
		return nil, errors.Errorf("GET is not 200 status: %d", postRes.StatusCode)
	}
	getUrl, err := util.UrlJoin(server, peerPath)
	if err != nil {
		return nil, err
	}
	getRes, err := http.Get(getUrl)
	if err != nil {
		return nil, err
	}
	if getRes.StatusCode != 200 {
		return nil, errors.Errorf("POST is not 200 status: %d", postRes.StatusCode)
	}
	return getRes.Body, nil
}
