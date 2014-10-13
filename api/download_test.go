package zippo

import (
	"gopkg.in/check.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

type DownloadSuite struct {
	server *httptest.Server
}

func init() {
	check.Suite(&DownloadSuite{})
}

func (s *DownloadSuite) SetUpSuite(c *check.C) {
	m := http.NewServeMux()
	m.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("fixtures/logo.png")
		if err != nil {
			panic(err)
		}

		defer f.Close()
		io.Copy(w, f)
	})

	s.server = httptest.NewServer(m)
}

func (s *DownloadSuite) TearDownSuite(c *check.C) {
	s.server.Close()
}

func (s *DownloadSuite) TestDownloadTemp(c *check.C) {
	p := Payload{
		Filename:    "logo.png",
		URL:         s.server.URL + "/logo.png",
		ContentType: "image/png",
	}

	f, err := DownloadTmp(p)
	c.Assert(err, check.IsNil)
	c.Assert(f, check.Not(check.Equals), "")
}
