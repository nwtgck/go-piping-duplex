package util

import (
	"golang.org/x/crypto/openpgp"
	"io"
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

func OpenpgpSymmetricallyEncrypt(plain io.Reader, passphrase []byte) io.Reader {
	// (base: https://gist.github.com/eliquious/9e96017f47d9bd43cdf9)
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		w, err := openpgp.SymmetricallyEncrypt(pw, passphrase, nil, nil)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(w, plain)
		if err != nil {
			panic(err)
		}
		w.Close()
	}()

	return pr
}

func OpenpgpSymmetricallyDecrypt(encrypted io.Reader, passphrase []byte) (io.Reader, error) {
	// (base: https://github.com/golang/crypto/blob/a2144134853fc9a27a7b1e3eb4f19f1a76df13c9/openpgp/write_test.go#L129)
	md, err := openpgp.ReadMessage(encrypted, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return passphrase, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	return md.UnverifiedBody, nil
}
