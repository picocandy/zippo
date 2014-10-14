package zippo

import (
	"gopkg.in/check.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type BaseSuite struct {
	server *httptest.Server
}

func (s *BaseSuite) SetUpSuite(c *check.C) {
	h := func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("fixtures/logo.png")
		if err != nil {
			panic(err)
		}

		defer f.Close()
		io.Copy(w, f)
	}

	n := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized!", http.StatusForbidden)
	}

	m := http.NewServeMux()
	m.HandleFunc("/logo.png", h)
	m.HandleFunc("/image.png", h)
	m.HandleFunc("/blocked.png", n)
	s.server = httptest.NewServer(m)
}

func (s *BaseSuite) TearDownSuite(c *check.C) {
	s.server.Close()
}

var fixtures = map[string]string{
	"archive": `
	{
		"filename": "zippo-archive.zip",
		"payloads": [
			{
				"url": "http://picocandy.com/images/logo.png",
				"filename": "picocandy.png",
				"content_type": "image/png"
			},
			{
				"url": "http://www.gorillatoolkit.org/static/images/gorilla-icon-64.png",
				"filename": "gorilla.png",
				"content_type": "image/png"
			}
		]
	}
	`,
	"archive-without-filename": `
	{
		"filename": "",
		"payloads": [
			{
				"url": "http://picocandy.com/images/logo.png",
				"filename": "picocandy.png",
				"content_type": "image/png"
			},
			{
				"url": "http://www.gorillatoolkit.org/static/images/gorilla-icon-64.png",
				"filename": "gorilla.png",
				"content_type": "image/png"
			}
		]
	}
	`,
}
