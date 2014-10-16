package zippo

import (
	"flag"
	"gopkg.in/check.v1"
	"io/ioutil"
	"testing"
)

var live = flag.Bool("live", false, "Include live tests")

func Test(t *testing.T) { check.TestingT(t) }

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

var hashString = "0e3e239fac217bf5396b4e670cbf5ac7ce7dface"
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
				"content_type": "image/png",
				"content_length": 3909
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
}
