package zippo

import (
	"encoding/json"
	"github.com/ncw/swift"
	"gopkg.in/check.v1"
	"net/http/httptest"
	"os"
)

type ArchiveSuite struct {
	server *httptest.Server
	cf     swift.Connection
}

func init() {
	check.Suite(&ArchiveSuite{cf: NewConnection()})
}

func (s *ArchiveSuite) TestArchive_String_withoutFilename(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive-without-filename"]), a)
	c.Assert(err, check.IsNil)
	c.Assert(a.String(), check.Equals, hashString+".zip")
}

func (s *ArchiveSuite) TestArchive_String(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive"]), a)
	c.Assert(err, check.IsNil)
	c.Assert(a.String(), check.Equals, "zippo-archive.zip")
}

func (s *ArchiveSuite) TestArchive_Hash(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive-without-filename"]), a)
	c.Assert(err, check.IsNil)
	c.Assert(a.Hash(), check.Equals, hashString)
}

func (s *ArchiveSuite) TestArchive_Build_failure(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive-failure"]), a)
	c.Assert(err, check.IsNil)

	err = a.Build()
	c.Assert(err, check.NotNil)
	c.Assert(a.TempFile, check.Equals, "")
}

func (s *ArchiveSuite) TestArchive_Build(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive"]), a)
	c.Assert(err, check.IsNil)

	err = a.Build()
	c.Assert(err, check.IsNil)
	c.Assert(a.TempFile, check.Matches, "(.)*zippo-archive.zip(.)*")

	n, err := os.Stat(a.TempFile)
	c.Assert(os.IsExist(err), check.Equals, false)
	c.Assert(n.Size(), check.Not(check.Equals), 22) // empty zip
}

func (s *ArchiveSuite) TestArchive_Upload(c *check.C) {
	if !*live {
		c.Skip("-live is not provided")
	}

	a := &Archive{}
	err := json.Unmarshal([]byte(fixtures["archive"]), a)
	c.Assert(err, check.IsNil)

	a.TempFile = prepareTemp("fixtures/zippo-archive.zip", "zippo-archive-suite-")
	a.SetConnection(NewConnection())
	a.Authenticate()

	o, h, err := a.Upload(container)
	c.Assert(err, check.IsNil)

	c.Assert(o.Name, check.Equals, "zippo-archive.zip")
	c.Assert(o.ContentType, check.Equals, "application/zip")
	c.Assert(o.Bytes, check.Equals, int64(10905))

	c.Assert(h.ObjectMetadata()["archive-hash"], check.Equals, hashString)
	c.Assert(h["Content-Type"], check.Equals, "application/zip")
	c.Assert(h["Content-Length"], check.Equals, "10905")
}

func (s *ArchiveSuite) TestArchive_RemoveTemp_failure(c *check.C) {
	a := &Archive{}
	err := a.RemoveTemp()
	c.Assert(err, check.NotNil)
}

func (s *ArchiveSuite) TestArchive_RemoveTemp(c *check.C) {
	t := prepareTemp("fixtures/logo.png", "zippo-archive-suite-")
	a := &Archive{}
	a.TempFile = t

	err := a.RemoveTemp()
	c.Assert(err, check.IsNil)
	c.Assert(a.TempFile, check.Equals, "")

	_, err = os.Stat(t)
	c.Assert(err, check.NotNil)
	c.Assert(os.IsNotExist(err), check.Equals, true)
}

func (s *ArchiveSuite) TestArchive_RenameDuplicatePayloads(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive-duplicate"]), a)
	c.Assert(err, check.IsNil)

	a.RenameDuplicatePayloads()

	c.Assert(a.Payloads[0].Filename, check.Equals, "picocandy.png")
	c.Assert(a.Payloads[1].Filename, check.Equals, "picocandy-1.png")
	c.Assert(a.Payloads[2].Filename, check.Equals, "picocandy.gif")
	c.Assert(a.Payloads[3].Filename, check.Equals, "Picocandy-2.png")
}
