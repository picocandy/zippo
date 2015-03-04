package zippo

import (
	"archive/zip"
	"bytes"
	"github.com/ncw/swift"
	"gopkg.in/check.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

type PayloadSuite struct {
	server *httptest.Server
	cf     swift.Connection
}

func init() {
	check.Suite(&PayloadSuite{cf: NewConnection()})
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

func (s *PayloadSuite) TestPayload_String_url(c *check.C) {
	p := &Payload{
		URL: "http://example.com/images/logo.png",
	}

	c.Assert(p.String(), check.Equals, "logo.png")
}

func (s *PayloadSuite) TestPayload_String(c *check.C) {
	p := &Payload{
		Filename: "picocandy_logo.png",
		URL:      "http://example.com/images/logo.png",
	}

	c.Assert(p.String(), check.Equals, "picocandy_logo.png")
}

func (s *PayloadSuite) TestPayload_Download_failure(c *check.C) {
	p := &Payload{
		Filename:    "unknown.png",
		URL:         s.server.URL + "/unknown.png",
		ContentType: "image/png",
	}

	err := p.Build()
	c.Assert(err, check.NotNil)
	c.Assert(p.TempFile, check.Equals, "")
}

func (s *PayloadSuite) TestPayload_Download_sizeMismatch(c *check.C) {
	p := &Payload{
		Filename:      "logo.png",
		URL:           s.server.URL + "/logo.png",
		ContentLength: 10,
	}

	err := p.Build()
	c.Assert(err, check.NotNil)
	c.Assert(p.TempFile, check.Equals, "")
}

func (s *PayloadSuite) TestPayload_Download_sizeAuto(c *check.C) {
	p := &Payload{
		Filename:      "logo.png",
		URL:           s.server.URL + "/logo.png",
		ContentLength: -1,
	}

	err := p.Build()
	c.Assert(err, check.IsNil)
	c.Assert(p.ContentLength, check.Equals, int64(139100))
	c.Assert(p.TempFile, check.Not(check.Equals), "")
}

func (s *PayloadSuite) TestPayload_Download(c *check.C) {
	p := &Payload{
		Filename:    "logo.png",
		URL:         s.server.URL + "/logo.png",
		ContentType: "image/png",
	}

	err := p.Build()
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
		TempFile: prepareTemp("fixtures/logo.png", "zippo-payload-suite-"),
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

	c.Assert(br.Len(), check.Not(check.Equals), 22) // empty zip

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
	t := prepareTemp("fixtures/logo.png", "zippo-payload-suite-")
	p := &Payload{TempFile: t}

	err := p.RemoveTemp()
	c.Assert(err, check.IsNil)
	c.Assert(p.TempFile, check.Equals, "")

	_, err = os.Stat(t)
	c.Assert(err, check.NotNil)
	c.Assert(os.IsNotExist(err), check.Equals, true)
}

func (s *PayloadSuite) TestPayload_Upload(c *check.C) {
	t := prepareTemp("fixtures/logo.png", "zippo-payload-suite-")
	p := &Payload{
		Filename:      "picocandy_logo.png",
		URL:           "http://picocandy.com/images/logo.png",
		ContentLength: 139100,
		TempFile:      t,
	}

	p.SetConnection(NewConnection())
	p.Authenticate()

	o, h, err := p.Upload(container)
	c.Assert(err, check.IsNil)

	c.Assert(o.Name, check.Equals, "picocandy_logo.png")
	c.Assert(o.ContentType, check.Equals, "image/png")
	c.Assert(o.Bytes, check.Equals, int64(139100))

	c.Assert(h.ObjectMetadata()["payload-hash"], check.Equals, "be296bc2ea9cf42eb3a292c387ffedb718959f69")
	c.Assert(h["Content-Type"], check.Equals, "image/png")
	c.Assert(h["Content-Length"], check.Equals, "139100")
}
