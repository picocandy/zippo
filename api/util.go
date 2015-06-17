package zippo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ncw/swift"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func GenerateTempURL(cf swift.Connection, r Parker) (string, error) {
	var err error

	u, err := url.Parse(cf.Auth.StorageUrl(false))
	if err != nil {
		return "", err
	}

	u.Path = fmt.Sprintf("%s/%s/%s", u.Path, container, r.String())

	method := "GET"
	expires := int(time.Now().Unix() + r.ExpirationSec())
	body := fmt.Sprintf("%s\n%d\n%s", method, expires, u.Path)

	h := hmac.New(sha1.New, []byte(metaTempKey))
	io.WriteString(h, body)

	v := url.Values{
		"temp_url_sig":     []string{hex.EncodeToString(h.Sum(nil))},
		"temp_url_expires": []string{strconv.Itoa(expires)},
	}

	u.RawQuery = v.Encode()
	return u.String(), nil
}

func NewConnection() swift.Connection {
	var internal bool
	if os.Getenv("SWIFT_INTERNAL") == "true" {
		internal = true
	} else {
		internal = false
	}

	return swift.Connection{
		UserName: os.Getenv("SWIFT_API_USER"),
		ApiKey:   os.Getenv("SWIFT_API_KEY"),
		AuthUrl:  os.Getenv("SWIFT_AUTH_URL"),
		Region:   os.Getenv("SWIFT_REGION"),
		TenantId: os.Getenv("SWIFT_TENANT_ID"),
		Internal: internal,
	}
}

func UpdateAccountMetaTempURL(cf swift.Connection) error {
	var err error

	if metaTempKey == "" {
		return errors.New("Missing SWIFT_META_TEMP value")
	}

	h := swift.Headers{"X-Account-Meta-Temp-Url-Key": metaTempKey}
	err = cf.AccountUpdate(h)
	if err != nil {
		return err
	}

	return nil
}

func JSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	s, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(s)
}

func PostJSON(u string, v interface{}) error {
	s, err := json.Marshal(v)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u, bytes.NewBuffer(s))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	// When err is nil, resp always contains a non-nil resp.Body.
	// We should close the resp.Body even if we don't read it
	defer resp.Body.Close()

	return nil
}

func SplitFilename(s string) (string, string) {
	file := s
	ext := ""

	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			file = s[:i]
			ext = s[i:]

			break
		}
	}

	return file, ext
}
