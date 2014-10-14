package zippo

import (
	"gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type HTTPSuite struct {
	server *httptest.Server
}

func init() {
	check.Suite(&HTTPSuite{})
}

func (s *HTTPSuite) SetUpSuite(c *check.C) {
	m := http.NewServeMux()
	m.HandleFunc("/", HomeHandler)
	m.HandleFunc("/z", func(w http.ResponseWriter, r *http.Request) {
		ZipHandler(w, r, NewConnection())
	})

	s.server = httptest.NewServer(m)
}

func (s *HTTPSuite) TearDownSuite(c *check.C) {
	s.server.Close()
}

func (s *HTTPSuite) TestHomeHandler(c *check.C) {
	resp, err := http.Get(s.server.URL)
	c.Assert(err, check.IsNil)
	defer resp.Body.Close()

	c.Assert(resp.StatusCode, check.Equals, http.StatusOK)

	b, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	c.Assert(string(b), check.Equals, "Zippo!")
}

func (s *HTTPSuite) TestZipHandler(c *check.C) {
	if !*live {
		c.Skip("-live is not provided")
	}

	req, err := http.NewRequest("POST", s.server.URL+"/z", strings.NewReader(fixtures["archive"]))
	c.Assert(err, check.IsNil)
	req.Header.Add("Content-Type", "application/json")

	hc := http.DefaultClient

	resp, err := hc.Do(req)
	c.Assert(err, check.IsNil)

	defer resp.Body.Close()

	c.Assert(resp.StatusCode, check.Equals, http.StatusOK)

	b, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	c.Assert(string(b), check.Matches, "(.)*zippo-archive.zip(.)*")
}
