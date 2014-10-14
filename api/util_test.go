package zippo

import (
	"encoding/json"
	"github.com/ncw/swift"
	"gopkg.in/check.v1"
	"net/url"
	"path"
)

type UtilSuite struct {
	cf swift.Connection
}

func init() {
	check.Suite(&UtilSuite{cf: NewConnection()})
}

func (s *UtilSuite) TestUtil_GenerateTempURL(c *check.C) {
	if !*live {
		c.Skip("-live is not provided")
	}

	a := &Archive{}

	err := json.Unmarshal([]byte(fixtures["archive-without-filename"]), a)
	c.Assert(err, check.IsNil)

	f, err := GenerateTempURL(s.cf, a)
	c.Assert(err, check.IsNil)

	u, err := url.Parse(f)
	c.Assert(err, check.IsNil)

	c.Assert(u.Query().Get("temp_url_sig"), check.Not(check.Equals), "")
	c.Assert(u.Query().Get("temp_url_expires"), check.Not(check.Equals), "")
	c.Assert(path.Base(u.Path), check.Equals, "24c6e8fcb0a625d23d2aff43ec487a90167d56bb.zip")
}
