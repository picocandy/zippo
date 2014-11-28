package zippo

import (
	"flag"
	"gopkg.in/check.v1"
	"io/ioutil"
	"testing"
)

var live = flag.Bool("live", false, "Include live tests")

func Test(t *testing.T) { check.TestingT(t) }

func init() {
	log.Out = ioutil.Discard
}

func prepareTemp(f string, prefix string) string {
	tmp, err := ioutil.TempFile("", prefix)
	if err != nil {
		panic(err)
	}
	defer tmp.Close()

	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	_, err = tmp.Write(b)
	if err != nil {
		panic(err)
	}

	return tmp.Name()
}

var hashString = "fe8f6f64250d93af797e1609c8839b6de7955967"
var fixtures = map[string]string{
	"archive": `
	{
		"filename": "zippo-archive.zip",
		"payloads": [
			{
				"url": "http://picocandy.com/images/logo.png",
				"filename": "picocandy.png",
				"content_type": "image/png",
				"content_length": 3909
			},
			{
				"url": "http://www.gorillatoolkit.org/static/images/gorilla-icon-64.png",
				"filename": "gorilla.png",
				"content_type": "image/png",
				"content_length": 6722
			}
		],
		"expiration": 60
	}
	`,
	"archive-without-filename": `
	{
		"filename": "",
		"payloads": [
			{
				"url": "http://picocandy.com/images/logo.png",
				"filename": "picocandy.png",
				"content_type": "image/png",
				"content_length": 3909
			},
			{
				"url": "http://www.gorillatoolkit.org/static/images/gorilla-icon-64.png",
				"filename": "gorilla.png",
				"content_type": "image/png",
				"content_length": 6722
			}
		],
		"expiration": 60
	}
	`,
	"archive-failure": `
	{
		"filename": "zippo-failure.zip",
		"payloads": [
			{
				"url": "http://picocandy.com/images/unknown.png",
				"filename": "picocandy.png",
				"content_type": "image/png"
			},
			{
				"url": "http://www.gorillatoolkit.org/static/images/gorilla-icon-64.png",
				"filename": "gorilla.png",
				"content_type": "image/png",
				"content_length": 6722
			}
		]
	}
	`,
	"payload": `
	{
		"url": "http://picocandy.com/images/logo.png",
		"filename": "picocandy.png",
		"content_type": "image/png",
		"content_length": 3909,
		"expiration": 60
	}
	`,
	"payload-without-filename": `
	{
		"url": "http://picocandy.com/images/logo.png",
		"content_type": "image/png",
		"content_length": 3909,
		"expiration": 60
	}
	`,
}
