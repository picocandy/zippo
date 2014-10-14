package zippo

import (
	"archive/zip"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Archive struct {
	Filename string    `json:"filename,omitempty"`
	Payloads []Payload `json:"payloads"`
	TempFile string    `json:"-"`
	Hash     string    `json:"-"`
}

func (a *Archive) String() string {
	if a.Filename != "" {
		return a.Filename
	}

	h := sha1.New()

	for _, p := range a.Payloads {
		io.WriteString(h, p.String())
	}

	a.Hash = fmt.Sprintf("%x.zip", h.Sum(nil))
	return a.Hash
}

func (a *Archive) Build() error {
	out, err := ioutil.TempFile("", a.String())
	if err != nil {
		return err
	}
	defer out.Close()

	z := zip.NewWriter(out)

	for _, p := range a.Payloads {
		err := p.Download()
		if err != nil {
			return err
		}

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