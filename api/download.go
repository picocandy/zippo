package zippo

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func DownloadTmp(pl Payload) (f string, err error) {
	out, err := ioutil.TempFile("", pl.Filename)
	if err != nil {
		return
	}

	defer out.Close()

	resp, err := http.Get(pl.URL)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to download %s, got %s", pl.URL, resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	return out.Name(), nil
}
