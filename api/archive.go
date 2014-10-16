package zippo

import (
	"archive/zip"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/ncw/swift"
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

	return a.SumHash() + ".zip"
}

func (a *Archive) SumHash() string {
	if a.Hash != "" {
		return a.Hash
	}

	h := sha1.New()

	for _, p := range a.Payloads {
		io.WriteString(h, p.String())
	}

	a.Hash = fmt.Sprintf("%x", h.Sum(nil))
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

func (a *Archive) Upload(cf swift.Connection, cn string) (ob swift.Object, h swift.Headers, err error) {
	f, err := os.Open(a.TempFile)
	if err != nil {
		return
	}
	defer f.Close()

	d := swift.Headers{"X-Object-Meta-Archive-Hash": a.SumHash()}
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

func (a *Archive) DownloadURL(cf swift.Connection) (string, error) {
	var err error

	i, h, err := cf.Object(container, a.String())
	if err != nil {
		return "", err
	}

	if i.Bytes == 0 || i.Bytes == 22 {
		return "", errors.New("Empty file detected")
	}

	if h.ObjectMetadata()["archive-hash"] != a.SumHash() {
		return "", errors.New("File is updated")
	}

	return GenerateTempURL(cf, a)
}
