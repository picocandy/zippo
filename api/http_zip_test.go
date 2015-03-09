package zippo

import (
	"encoding/json"
	"fmt"
	"gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *HTTPSuite) TestHandler_ZipUpload(c *check.C) {
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

func (s *HTTPSuite) TestHandler_ZipUpload_callback(c *check.C) {
	if !*live {
		c.Skip("-live is not provided")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t Response

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&t)
		c.Assert(err, check.IsNil)

		fmt.Fprintln(w, "Hello Callback!")
	}))
	defer ts.Close()

	data := fixtures["archive-with-callback"]
	data = strings.Replace(data, "http://example.com/", ts.URL, 1)

	req, err := http.NewRequest("POST", s.server.URL+"/z", strings.NewReader(data))
	c.Assert(err, check.IsNil)
	req.Header.Add("Content-Type", "application/json")

	hc := http.DefaultClient

	resp, err := hc.Do(req)
	c.Assert(err, check.IsNil)

	defer resp.Body.Close()

	c.Assert(resp.StatusCode, check.Equals, http.StatusAccepted)

	b, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, check.IsNil)
	c.Assert(string(b), check.Equals, `{"status":202,"message":"Request is being processed."}`)
}
