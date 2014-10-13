package zippo

import (
	"gopkg.in/check.v1"
)

type DownloadSuite struct {
	suite *BaseSuite
}

func init() {
	check.Suite(&DownloadSuite{suite: &BaseSuite{}})
}

func (s *DownloadSuite) SetUpSuite(c *check.C) {
	s.suite.SetUpSuite(c)
}

func (s *DownloadSuite) TearDownSuite(c *check.C) {
	s.suite.TearDownSuite(c)
}

func (s *DownloadSuite) TestDownloadTemp_failure(c *check.C) {
	p := &Payload{
		Filename:    "unknown.png",
		URL:         s.suite.server.URL + "/unknown.png",
		ContentType: "image/png",
	}

	err := DownloadTmp(p)
	c.Assert(err, check.NotNil)
	c.Assert(p.TempFile, check.Equals, "")
}

func (s *DownloadSuite) TestDownloadTemp(c *check.C) {
	p := &Payload{
		Filename:    "logo.png",
		URL:         s.suite.server.URL + "/logo.png",
		ContentType: "image/png",
	}

	err := DownloadTmp(p)
	c.Assert(err, check.IsNil)
	c.Assert(p.TempFile, check.Not(check.Equals), "")
}
