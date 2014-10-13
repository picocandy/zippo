package zippo

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func ZipBuilder(ps []*Payload) (string, error) {
	h := zipHash(ps)
	out, err := ioutil.TempFile("", h)
	if err != nil {
		return "", err
	}

	defer out.Close()

	z := zip.NewWriter(out)

	for i := range ps {
		p := ps[i]

		err := DownloadTmp(p)
		if err != nil {
			return "", err
		}

		f, err := z.Create(p.Filename)
		if err != nil {
			return "", err
		}

		ot, err := os.Open(p.TempFile)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(f, ot)
		if err != nil {
			return "", err
		}

		if err = ot.Close(); err == nil {
			p.RemoveTemp()
		}
	}

	if err = z.Close(); err != nil {
		return "", err
	}

	return out.Name(), nil
}

func zipHash(ps []*Payload) string {
	h := sha1.New()

	for i := range ps {
		io.WriteString(h, ps[i].String())
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
