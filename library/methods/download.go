package methods

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func DownloadImg(name, src, path string) error {
	proxy, _ := url.Parse("http://127.0.0.1:9999")
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
		Timeout: time.Minute,
	}
	resp, err := c.Get(src)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(`http get err name: %s; url: %s`, name, src))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(`read body err name: %s; url: %s`, name, src))

	}

	idx := strings.LastIndex(src, "/")
	var ext string
	spilt := strings.Split(src[idx:], ".")
	if len(spilt) > 0 {
		ext = spilt[1]
	}

	imgName := fmt.Sprintf(`%s%s.%s`, path, name, ext)
	fmt.Println(imgName)
	if err = ioutil.WriteFile(imgName, body, 0666); err != nil {
		return errors.Wrap(err, fmt.Sprintf(`download err name: %s; url: %s`, name, src))

	}

	return nil
}

func DownloadBook(path string, url string) {

}
