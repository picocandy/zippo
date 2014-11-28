package zippo

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/ncw/swift"
	"io"
	"io/ioutil"
	"os"
)

type Archive struct {
	Filename   string     `json:"filename,omitempty"`
	Payloads   []*Payload `json:"payloads"`
	TempFile   string     `json:"-"`
	Expiration int64      `json:"expiration,omitempty"`
	hash       string
}

func (a *Archive) String() string {
	if a.Filename != "" {
		return a.Filename
	}

	return a.Hash() + ".zip"
}

func (a *Archive) Hash() string {
	if a.hash != "" {
		return a.hash
	}

	h := sha1.New()

	for _, p := range a.Payloads {
		io.WriteString(h, p.Hash())
	}

	a.hash = hex.EncodeToString(h.Sum(nil))
	return a.hash
}

func (a *Archive) Build() error {
	out, err := ioutil.TempFile("", a.String())
	if err != nil {
		return err
	}
	defer out.Close()

	z := zip.NewWriter(out)

	c := make(chan error)

	for _, p := range a.Payloads {
		go func(p *Payload) {
			c <- p.Download()
		}(p)
	}

	for i := 1; i <= len(a.Payloads); i++ {
		err := <-c
		if err != nil {
			return err
		}
	}

	for _, p := range a.Payloads {
		err = p.WriteZip(z)
		if err != nil {
			return err
		}

		err = p.RemoveTemp()
		if err != nil {
			return err
		}
	}

	if err = z.Close(); err != nil {
		return err
	}

	a.TempFile = out.Name()
	return nil
}

func (a *Archive) Upload(cf swift.Connection, cn string) (ob swift.Object, h swift.Headers, err error) {
	f, err := os.Open(a.TempFile)
	if err != nil {
		return
	}
	defer f.Close()

	d := swift.Headers{"X-Object-Meta-Archive-Hash": a.Hash()}
	_, err = cf.ObjectPut(cn, a.String(), f, true, "", "application/zip", d)
	if err != nil {
		return
	}

	return cf.Object(cn, a.String())
}

func (a *Archive) RemoveTemp() error {
	if a.TempFile == "" {
		return errors.New("No valid temporary file available")
	}

	err := os.Remove(a.TempFile)
	if err == nil {
		a.TempFile = ""
	}

	return err
}

func (a *Archive) ExpirationSec() int64 {
	if a.Expiration == 0 {
		return 600
	}

	return a.Expiration
}

func (a *Archive) DownloadURL(cf swift.Connection) (string, error) {
	var err error

	i, h, err := cf.Object(container, a.String())
	if err != nil {
		return "", err
	}

	if i.Bytes == 0 || i.Bytes == 22 {
		return "", errors.New("Empty file detected")
	}

	if h.ObjectMetadata()["archive-hash"] != a.Hash() {
		return "", errors.New("File is updated")
	}

	return GenerateTempURL(cf, a)
}
