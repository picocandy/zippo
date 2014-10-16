package zippo

import (
	"encoding/json"
	"github.com/ncw/swift"
	"gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
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
	c.Assert(path.Base(u.Path), check.Equals, hashString+".zip")
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

func (s *UtilSuite) TestJSON(c *check.C) {
	type Notice struct {
		Message string `json:"message"`
	}

	f := map[int]interface{}{
		http.StatusOK:           Notice{Message: "well done"},
		http.StatusUnauthorized: map[string]string{"error": "unauthorized"},
	}

	e := map[int]string{
		http.StatusOK:           `{"message":"well done"}`,
		http.StatusUnauthorized: `{"error":"unauthorized"}`,
	}

	for k, v := range f {
		w := httptest.NewRecorder()
		JSON(w, v, k)
		c.Assert(w.Code, check.Equals, k)
		c.Assert(w.Header().Get("Content-Type"), check.Equals, "application/json; charset=UTF-8")
		c.Assert(w.Body.String(), check.Equals, e[k])
	}
}
