package zippo

import (
	"errors"
	"github.com/ncw/swift"
	"io"
	"io/ioutil"
	"os"
)

type Expiration struct {
	Duration int64 `json:"expiration,omitempty"`
}

func (e *Expiration) ExpirationSec() int64 {
	if e.Duration == 0 {
		return 600
	}

	return e.Duration
}

type Callback struct {
	CallbackURL string `json:"callback_url,omitempty"`
}

func (c *Callback) HasCallbackURL() bool {
	return c.CallbackURL != ""
}

func (c *Callback) CallCallbackURL(data Response) error {
	return PostJSON(c.CallbackURL, data)
}

type CloudFile struct {
	cf        swift.Connection
	container string
}

func (c *CloudFile) SetConnection(conn swift.Connection) {
	c.cf = conn
}

func (c *CloudFile) SetContainer(str string) {
	c.container = str
}

func (c *CloudFile) Authenticate() error {
	return c.cf.Authenticate()
}

func (c *CloudFile) Container() string {
	return c.container
}

type Temporary struct {
	TempFile string `json:"-"`
}

func (t *Temporary) WriteTemp(str string, data io.Reader) error {
	out, err := t.CreateTemp(str)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, data)
	if err != nil {
		return err
	}

	return nil
}

func (t *Temporary) CreateTemp(str string) (*os.File, error) {
	out, err := ioutil.TempFile("", str)
	if err != nil {
		return nil, err
	}

	t.TempFile = out.Name()
	return out, nil
}

func (t *Temporary) TempStat() (os.FileInfo, error) {
	return os.Stat(t.TempFile)
}

func (t *Temporary) RemoveTemp() error {
	if t.TempFile == "" {
		return errors.New("No valid temporary file available")
	}

	err := os.Remove(t.TempFile)
	if err == nil {
		t.TempFile = ""
	}

	return err
}
