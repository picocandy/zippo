package zippo

import (
	"encoding/json"
	"gopkg.in/check.v1"
	"net/http/httptest"
)

type ArchiveSuite struct {
	server *httptest.Server
}

func init() {
	check.Suite(&ArchiveSuite{})
}

func (s *ArchiveSuite) TestString_withoutFilename(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive-without-filename"]), a)
	c.Assert(err, check.IsNil)
	c.Assert(a.String(), check.Equals, "24c6e8fcb0a625d23d2aff43ec487a90167d56bb.zip")
}

func (s *ArchiveSuite) TestString(c *check.C) {
	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive"]), a)
	c.Assert(err, check.IsNil)
	c.Assert(a.String(), check.Equals, "zippo-archive.zip")
}
