package zippo

import (
	"archive/zip"
	"bytes"
	"gopkg.in/check.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

type PayloadSuite struct {
	server *httptest.Server
}

func init() {
	check.Suite(&PayloadSuite{})
}

func (s *PayloadSuite) SetUpSuite(c *check.C) {
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

func (s *PayloadSuite) TearDownSuite(c *check.C) {
	s.server.Close()
}

func (s *PayloadSuite) TestPayload_String(c *check.C) {
	p := &Payload{
		Filename: "picocandy_logo.png",
		URL:      "http://example.com/images/logo.png",
	}

	want := "picocandy_logo.png::http://example.com/images/logo.png"
	c.Assert(p.String(), check.Equals, want)
}

func (s *PayloadSuite) TestPayload_Download_failure(c *check.C) {
	p := &Payload{
		Filename:    "unknown.png",
		URL:         s.server.URL + "/unknown.png",
		ContentType: "image/png",
	}

	err := p.Download()
	c.Assert(err, check.NotNil)
	c.Assert(p.TempFile, check.Equals, "")
}

func (s *PayloadSuite) TestPayload_Download(c *check.C) {
	p := &Payload{
		Filename:    "logo.png",
		URL:         s.server.URL + "/logo.png",
		ContentType: "image/png",
	}

	err := p.Download()
	c.Assert(err, check.IsNil)
	c.Assert(p.TempFile, check.Not(check.Equals), "")
}

func (s *PayloadSuite) TestPayload_WriteZip_failure(c *check.C) {
	p := &Payload{Filename: "awesome-logo.png"}

	buf := new(bytes.Buffer)
	z := zip.NewWriter(buf)

	err := p.WriteZip(z)
	c.Assert(err, check.NotNil)

	err = z.Close()
	if err != nil {
		panic(err)
	}

	c.Assert(buf.Len(), check.Equals, 22) // empty zip
}

func (s *PayloadSuite) TestPayload_WriteZip(c *check.C) {
	p := &Payload{
		Filename: "awesome-logo.png",
		TempFile: prepareTemp("zippo-payload-suite-"),
	}

	buf := new(bytes.Buffer)
	z := zip.NewWriter(buf)

	err := p.WriteZip(z)
	c.Assert(err, check.IsNil)

	err = z.Close()
	if err != nil {
		panic(err)
	}

	// check zip integrity
	br := bytes.NewReader(buf.Bytes())

	c.Assert(br.Len(), check.Not(check.Equals), 22)

	oz, err := zip.NewReader(br, int64(br.Len()))
	c.Assert(err, check.IsNil)

	c.Assert(len(oz.File), check.Equals, 1)
	c.Assert(oz.File[0].Name, check.Equals, "awesome-logo.png")
}

func (s *PayloadSuite) TestPayload_RemoveTemp_failure(c *check.C) {
	p := &Payload{}
	err := p.RemoveTemp()
	c.Assert(err, check.NotNil)
}

func (s *PayloadSuite) TestPayload_RemoveTemp(c *check.C) {
	t := prepareTemp("zippo-payload-suite-")
	p := &Payload{TempFile: t}

	err := p.RemoveTemp()
	c.Assert(err, check.IsNil)
	c.Assert(p.TempFile, check.Equals, "")

	_, err = os.Stat(t)
	c.Assert(err, check.NotNil)
	c.Assert(os.IsNotExist(err), check.Equals, true)
}
