package zippo

import (
	"gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type HTTPSuite struct {
	server *httptest.Server
}

func init() {
	check.Suite(&HTTPSuite{})
}

func (s *HTTPSuite) SetUpSuite(c *check.C) {
	h := NewHandler(NewConnection())
	m := http.NewServeMux()
	m.HandleFunc("/", HomeHandler)
	m.HandleFunc("/z", h.ZipUpload)
	m.HandleFunc("/u", h.Upload)

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
	c.Assert(string(b), check.Equals, "zippo!")
}
