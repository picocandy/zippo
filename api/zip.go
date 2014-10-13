package zippo

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func ZipBuilder(ps []Payload) error {
	var err error

	h := zipHash(ps)
	out, err := ioutil.TempFile("", h)
	if err != nil {
		return err
	}

	defer out.Close()

	z := zip.NewWriter(out)

	for i := range ps {
		p := ps[i]

		t, _ := DownloadTmp(p)
		if err != nil {
			return err
		}

		f, err := z.Create(p.Filename)
		if err != nil {
			return err
		}

		ot, err := os.Open(t)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, ot)
		if err != nil {
			return err
		}

		ot.Close()
	}

	if err = z.Close(); err != nil {
		return err
	}

	return nil
}

func zipHash(ps []Payload) string {
	h := sha1.New()

	for i := range ps {
		io.WriteString(h, ps[i].String())
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
