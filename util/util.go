package util

import (
	"net/url"
	"path"
)

// (from: https://stackoverflow.com/a/34668130/2885946)
func UrlJoin(rawurl string, elem ...string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(append([]string{u.Path}, elem...)...)
	return u.String(), nil
}
