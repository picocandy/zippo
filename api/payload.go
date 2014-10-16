package zippo

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ncw/swift"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Payload struct {
	URL           string `json:"url"`
	Filename      string `json:"filename"`
	ContentType   string `json:"content_type"`
	ContentLength int    `json:"content_length,omitempty"`
	TempFile      string `json:"-"`
	hash          string
}

func (p *Payload) String() string {
	if p.Filename != "" {
		return p.Filename
	}

	return path.Base(p.URL)
}

func (p *Payload) Hash() string {
	if p.hash != "" {
		return p.hash
	}

	h := sha1.New()

	io.WriteString(h, p.Filename)
	io.WriteString(h, p.URL)
	io.WriteString(h, strconv.Itoa(p.ContentLength))

	p.hash = hex.EncodeToString(h.Sum(nil))
	return p.hash
}

func (p *Payload) Download() error {
	out, err := ioutil.TempFile("", p.Filename)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(p.URL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to download %s, got %s", p.URL, resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	p.TempFile = out.Name()
	return nil
}

func (p *Payload) WriteZip(z *zip.Writer) error {
	if p.TempFile == "" {
		return errors.New("No valid temporary file available")
	}

	f, err := z.Create(p.Filename)
	if err != nil {
		return err
	}

	t, err := os.Open(p.TempFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, t)
	if err != nil {
		return err
	}

	return t.Close()
}

func (p *Payload) RemoveTemp() error {
	if p.TempFile == "" {
		return errors.New("No valid temporary file available")
	}

	err := os.Remove(p.TempFile)
	if err == nil {
		p.TempFile = ""
	}

	return err
}

func (p *Payload) Upload(cf swift.Connection, cn string) (ob swift.Object, h swift.Headers, err error) {
	f, err := os.Open(p.TempFile)
	if err != nil {
		return
	}
	defer f.Close()

	d := swift.Headers{"X-Object-Meta-Payload-Hash": p.Hash()}
	_, err = cf.ObjectPut(cn, p.String(), f, true, "", p.ContentType, d)
	if err != nil {
		return
	}

	return cf.Object(cn, p.String())
}
