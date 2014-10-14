package zippo

import (
	"gopkg.in/check.v1"
	"io/ioutil"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

func prepareTemp(prefix string) string {
	tmp, err := ioutil.TempFile("", prefix)
	if err != nil {
		panic(err)
	}
	defer tmp.Close()

	b, err := ioutil.ReadFile("fixtures/logo.png")
	if err != nil {
		panic(err)
	}

	_, err = tmp.Write(b)
	if err != nil {
		panic(err)
	}

	return tmp.Name()
}

var fixtures = map[string]string{
	"archive": `
	{
		"filename": "zippo-archive.zip",
		"payloads": [
			{
				"url": "http://picocandy.com/images/logo.png",
				"filename": "picocandy.png",
				"content_type": "image/png"
			},
			{
				"url": "http://www.gorillatoolkit.org/static/images/gorilla-icon-64.png",
				"filename": "gorilla.png",
				"content_type": "image/png"
			}
		]
	}
	`,
	"archive-without-filename": `
	{
		"filename": "",
		"payloads": [
			{
				"url": "http://picocandy.com/images/logo.png",
				"filename": "picocandy.png",
				"content_type": "image/png"
			},
			{
				"url": "http://www.gorillatoolkit.org/static/images/gorilla-icon-64.png",
				"filename": "gorilla.png",
				"content_type": "image/png"
			}
		]
	}
	`,
}
