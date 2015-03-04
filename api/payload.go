package zippo

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ncw/swift"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Payload struct {
	Expiration
	Callback
	CloudFile
	Temporary
	URL           string `json:"url"`
	Filename      string `json:"filename"`
	ContentType   string `json:"content_type"`
	ContentLength int64  `json:"content_length,omitempty"`
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
	io.WriteString(h, strconv.FormatInt(p.ContentLength, 10))

	p.hash = hex.EncodeToString(h.Sum(nil))
	return p.hash
}

func (p *Payload) Build() error {
	resp, err := http.Get(p.URL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to download %s, got %s", p.URL, resp.Status)
	}

	err = p.WriteTemp(p.Filename, resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to save %s to %s: %s", p.URL, p.TempFile, err.Error())
	}

	if p.ContentLength != 0 {
		i, err := p.TempStat()
		if err != nil {
			return err
		}

		if p.ContentLength == -1 {
			p.ContentLength = i.Size()
		} else {
			if i.Size() != p.ContentLength {
				p.RemoveTemp()
				return fmt.Errorf("Content length mismatch, expected %d, got %d", p.ContentLength, i.Size())
			}
		}
	}

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

func (p *Payload) Upload(cn string) (ob swift.Object, h swift.Headers, err error) {
	f, err := os.Open(p.TempFile)
	if err != nil {
		return
	}
	defer f.Close()

	d := swift.Headers{"X-Object-Meta-Payload-Hash": p.Hash()}
	_, err = p.cf.ObjectPut(cn, p.String(), f, true, "", p.ContentType, d)
	if err != nil {
		return
	}

	return p.cf.Object(cn, p.String())
}

func (p *Payload) DownloadURL() (string, error) {
	var err error

	i, h, err := p.cf.Object(container, p.String())
	if err != nil {
		return "", err
	}

	if i.Bytes == 0 {
		return "", errors.New("Empty file detected")
	}

	if h.ObjectMetadata()["payload-hash"] != p.Hash() {
		return "", errors.New("File is updated")
	}

	return GenerateTempURL(p.cf, p)
}

func (p *Payload) LogFields() logrus.Fields {
	return logrus.Fields{
		"hash":           p.Hash(),
		"filename":       p.String(),
		"url":            p.URL,
		"content_type":   p.ContentType,
		"content_length": p.ContentLength,
		"expiration":     p.ExpirationSec(),
	}
}
