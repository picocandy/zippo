package zippo

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Payload struct {
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	TempFile    string `json:"-"`
}

func (p *Payload) String() string {
	return p.Filename + "::" + p.URL
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
