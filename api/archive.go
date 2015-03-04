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
	"os"
	"strings"
)

type Archive struct {
	Expiration
	Callback
	CloudFile
	Temporary
	Filename string     `json:"filename,omitempty"`
	Payloads []*Payload `json:"payloads"`
	hash     string
}

func NewArchive(cf swift.Connection) *Archive {
	return &Archive{CloudFile: CloudFile{cf: cf, container: container}}
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
	out, err := a.CreateTemp(a.String())
	if err != nil {
		return err
	}
	defer out.Close()

	z := zip.NewWriter(out)

	c := make(chan error)

	for _, p := range a.Payloads {
		go func(p *Payload) {
			p.CallbackURL = ""
			p.SetConnection(a.cf)
			c <- p.Build()
		}(p)
	}

	for i := 1; i <= len(a.Payloads); i++ {
		err := <-c
		if err != nil {
			a.RemoveTemp()
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
		return a.RemoveTemp()
	}

	return nil
}

func (a *Archive) Upload() (ob swift.Object, h swift.Headers, err error) {
	f, err := os.Open(a.TempFile)
	if err != nil {
		return
	}
	defer f.Close()

	d := swift.Headers{"X-Object-Meta-Archive-Hash": a.Hash()}
	_, err = a.cf.ObjectPut(a.Container(), a.String(), f, true, "", "application/zip", d)
	if err != nil {
		return
	}

	return a.cf.Object(a.Container(), a.String())
}

func (a *Archive) DownloadURL() (string, error) {
	var err error

	i, h, err := a.cf.Object(container, a.String())
	if err != nil {
		return "", err
	}

	if i.Bytes == 0 || i.Bytes == 22 {
		return "", errors.New("Empty file detected")
	}

	if h.ObjectMetadata()["archive-hash"] != a.Hash() {
		return "", errors.New("File is updated")
	}

	return GenerateTempURL(a.cf, a)
}

func (a *Archive) RenameDuplicatePayloads() {
	filenames := make(map[string]int)

	for _, p := range a.Payloads {
		key := strings.ToLower(p.Filename)

		_, present := filenames[key]
		if present {
			filenames[key]++
			file, ext := SplitFilename(p.Filename)
			p.Filename = fmt.Sprintf("%s-%d%s", file, filenames[key], ext)
		} else {
			filenames[key] = 0
		}
	}
}

func (a *Archive) LogFields() logrus.Fields {
	return logrus.Fields{
		"hash":         a.Hash(),
		"filename":     a.String(),
		"content_type": "application/zip",
		"expiration":   a.ExpirationSec(),
	}
}
