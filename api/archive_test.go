package zippo

import (
	"encoding/json"
	"gopkg.in/check.v1"
	"net/http/httptest"
	"os"
)

type ArchiveSuite struct {
	server *httptest.Server
}

func init() {
	check.Suite(&ArchiveSuite{})
}

func (s *ArchiveSuite) TestArchive_String_withoutFilename(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive-without-filename"]), a)
	c.Assert(err, check.IsNil)
	c.Assert(a.String(), check.Equals, "24c6e8fcb0a625d23d2aff43ec487a90167d56bb.zip")
}

func (s *ArchiveSuite) TestArchive_String(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive"]), a)
	c.Assert(err, check.IsNil)
	c.Assert(a.String(), check.Equals, "zippo-archive.zip")
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

func (s *ArchiveSuite) TestArchive_RemoveTemp_failure(c *check.C) {
	a := &Archive{}
	err := a.RemoveTemp()
	c.Assert(err, check.NotNil)
}

func (s *ArchiveSuite) TestArchive_RemoveTemp(c *check.C) {
	t := prepareTemp("zippo-archive-suite-")
	a := &Archive{TempFile: t}

	err := a.RemoveTemp()
	c.Assert(err, check.IsNil)
	c.Assert(a.TempFile, check.Equals, "")

	_, err = os.Stat(t)
	c.Assert(err, check.NotNil)
	c.Assert(os.IsNotExist(err), check.Equals, true)
}
