package zippo

import (
	"gopkg.in/check.v1"
)

type ZipSuite struct {
	suite *BaseSuite
}

func init() {
	check.Suite(&ZipSuite{suite: &BaseSuite{}})
}

func (s *ZipSuite) SetUpSuite(c *check.C) {
	s.suite.SetUpSuite(c)
}

func (s *ZipSuite) TearDownSuite(c *check.C) {
	s.suite.TearDownSuite(c)
}

func (s *ZipSuite) TestZipBuilder_failure(c *check.C) {
	ps := []Payload{
		Payload{
			Filename:    "logo.png",
			URL:         s.suite.server.URL + "/logo.png",
			ContentType: "image/png",
		},
		Payload{
			Filename:    "blocked.png",
			URL:         s.suite.server.URL + "/blocked.png",
			ContentType: "image/png",
		},
	}

	f, err := ZipBuilder(ps)
	c.Assert(err, check.NotNil)
	c.Assert(f, check.Equals, "")
}

func (s *ZipSuite) TestZipBuilder(c *check.C) {
	ps := []Payload{
		Payload{
			Filename:    "logo.png",
			URL:         s.suite.server.URL + "/logo.png",
			ContentType: "image/png",
		},
		Payload{
			Filename:    "image.png",
			URL:         s.suite.server.URL + "/image.png",
			ContentType: "image/png",
		},
	}

	f, err := ZipBuilder(ps)
	c.Assert(err, check.IsNil)
	c.Assert(f, check.Not(check.Equals), "")
}
