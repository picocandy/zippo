package zippo

import (
	"encoding/json"
	"github.com/ncw/swift"
	"gopkg.in/check.v1"
	"net/url"
	"os"
	"path"
)

type UtilSuite struct {
	cf swift.Connection
}

func init() {
	check.Suite(&UtilSuite{cf: NewConnection()})
}

func (s *UtilSuite) SetUpSuite(c *check.C) {
	err := s.cf.Authenticate()
	c.Assert(err, check.IsNil)
}

func (s *UtilSuite) TestGenerateTempURL(c *check.C) {
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

func (s *UtilSuite) TestUpdateAccountMetaTempURL(c *check.C) {
	if !*live {
		c.Skip("-live is not provided")
	}

	err := UpdateAccountMetaTempURL(s.cf)
	c.Assert(err, check.IsNil)

	_, h, err := s.cf.Account()
	c.Assert(err, check.IsNil)

	key := os.Getenv("SWIFT_META_TEMP")
	c.Assert(h.AccountMetadata()["temp-url-key"], check.Equals, key)
}
