package zippo

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/ncw/swift"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"
)

func GenerateTempURL(cf swift.Connection, a *Archive) (string, error) {
	var err error

	err = cf.Authenticate()
	if err != nil {
		return "", err
	}

	u, err := url.Parse(cf.Auth.StorageUrl(false))
	if err != nil {
		return "", err
	}

	container := os.Getenv("SWIFT_CONTAINER")
	u.Path = fmt.Sprintf("%s/%s/%s", u.Path, container, a.String())

	key := os.Getenv("SWIFT_META_TEMP")
	method := "GET"
	expires := int(time.Now().Unix() + 600)
	body := fmt.Sprintf("%s\n%d\n%s", method, expires, u.Path)

	h := hmac.New(sha1.New, []byte(key))
	io.WriteString(h, body)

	v := url.Values{
		"temp_url_sig":     []string{fmt.Sprintf("%x", h.Sum(nil))},
		"temp_url_expires": []string{strconv.Itoa(expires)},
	}

	u.RawQuery = v.Encode()
	return u.String(), nil
}

func NewConnection() swift.Connection {
	return swift.Connection{
		UserName: os.Getenv("SWIFT_API_USER"),
		ApiKey:   os.Getenv("SWIFT_API_KEY"),
		AuthUrl:  os.Getenv("SWIFT_AUTH_URL"),
		Region:   os.Getenv("SWIFT_REGION"),
		TenantId: os.Getenv("SWIFT_TENANT_ID"),
	}
}

func UpdateAccountMetaTempURL(cf swift.Connection) error {
	var err error

	err = cf.Authenticate()
	if err != nil {
		return err
	}

	key := os.Getenv("SWIFT_META_TEMP")
	if key == "" {
		return errors.New("Missing SWIFT_META_TEMP value")
	}

	h := swift.Headers{"X-Account-Meta-Temp-Url-Key": key}
	err = cf.AccountUpdate(h)
	if err != nil {
		return err
	}

	return nil
}
