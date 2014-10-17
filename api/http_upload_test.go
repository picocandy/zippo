package zippo

import (
	"gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s *HTTPSuite) TestUploadHandler(c *check.C) {
	if !*live {
		c.Skip("-live is not provided")
	}

	req, err := http.NewRequest("POST", s.server.URL+"/u", strings.NewReader(fixtures["payload"]))
	c.Assert(err, check.IsNil)
	req.Header.Add("Content-Type", "application/json")

	hc := http.DefaultClient

	resp, err := hc.Do(req)
	c.Assert(err, check.IsNil)

	defer resp.Body.Close()

	c.Assert(resp.StatusCode, check.Equals, http.StatusOK)

	b, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	c.Assert(string(b), check.Matches, "(.)*picocandy.png(.)*")
}
