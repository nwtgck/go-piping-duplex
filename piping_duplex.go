package piping_duplex

import (
	"github.com/nwtgck/go-piping-duplex/util"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

func DuplexReader(server string, uploadPath string, downloadPath string, r io.Reader) (io.Reader, <-chan error, error) {
	uploadFinishErrCh := make(chan error)
	postUrl, err := util.UrlJoin(server, uploadPath)
	if err != nil {
		return nil, uploadFinishErrCh, err
	}
	// TODO: hard code
	contentType := "application/octet-stream"
	postRes, err := http.Post(postUrl, contentType, r)
	if err != nil {
		return nil, uploadFinishErrCh, err
	}
	if postRes.StatusCode != 200 {
		return nil, uploadFinishErrCh, errors.Errorf("GET is not 200 status: %d", postRes.StatusCode)
	}
	go func() {
		_, err := io.Copy(ioutil.Discard, postRes.Body)
		uploadFinishErrCh <- err
	}()
	getUrl, err := util.UrlJoin(server, downloadPath)
	if err != nil {
		return nil, uploadFinishErrCh, err
	}
	getRes, err := http.Get(getUrl)
	if err != nil {
		return nil, uploadFinishErrCh, err
	}
	if getRes.StatusCode != 200 {
		return nil, uploadFinishErrCh, errors.Errorf("POST is not 200 status: %d", postRes.StatusCode)
	}
	return getRes.Body, uploadFinishErrCh, nil
}
