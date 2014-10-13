package zippo

import (
	"archive/zip"
	"bytes"
	"gopkg.in/check.v1"
	"io/ioutil"
	"os"
)

type PayloadSuite struct {
	suite *BaseSuite
}

func init() {
	check.Suite(&PayloadSuite{suite: &BaseSuite{}})
}

func (s *PayloadSuite) SetUpSuite(c *check.C) {
	s.suite.SetUpSuite(c)
}

func (s *PayloadSuite) TearDownSuite(c *check.C) {
	s.suite.TearDownSuite(c)
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
		URL:         s.suite.server.URL + "/unknown.png",
		ContentType: "image/png",
	}

	err := p.Download()
	c.Assert(err, check.NotNil)
	c.Assert(p.TempFile, check.Equals, "")
}

func (s *PayloadSuite) TestPayload_Download(c *check.C) {
	p := &Payload{
		Filename:    "logo.png",
		URL:         s.suite.server.URL + "/logo.png",
		ContentType: "image/png",
	}

	err := p.Download()
	c.Assert(err, check.IsNil)
	c.Assert(p.TempFile, check.Not(check.Equals), "")
}

func (s *PayloadSuite) TestPayload_WriteZip(c *check.C) {
	p := &Payload{
		Filename: "awesome-logo.png",
		TempFile: s.prepareTemp(),
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

	c.Assert(br.Len(), check.Not(check.Equals), 0)

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
	t := s.prepareTemp()
	p := &Payload{TempFile: t}

	err := p.RemoveTemp()
	c.Assert(err, check.IsNil)
	c.Assert(p.TempFile, check.Equals, "")

	_, err = os.Stat(t)
	c.Assert(err, check.NotNil)
	c.Assert(os.IsNotExist(err), check.Equals, true)
}

func (s *PayloadSuite) prepareTemp() string {
	tmp, err := ioutil.TempFile("", "zippo-payload-suite-")
	if err != nil {
		panic(err)
	}
	defer tmp.Close()

	b, err := ioutil.ReadFile("fixtures/logo.png")
	if err != nil {
		panic(err)
	}

	_, err = tmp.Write(b)
	if err != nil {
		panic(err)
	}

	return tmp.Name()
}
