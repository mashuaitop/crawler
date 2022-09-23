package methods

import (
	"context"
	"crawler/utils"
	"fmt"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
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

func DownloadBook(uid, key, url, dir string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36 Aoyou/cXRsNCdsM3s-T1c8SHhARZqOZNMwOHWB7sPpE_x2ULIWqtc__h71MI7ASQ==`),
		//chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if err := chromedp.Run(ctx,
		utils.Setcookies(".u1lib.org", "remix_userid", uid, "remix_userkey", key),
		chromedp.Navigate(url)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("打开书籍详情页失败: url: %s", url))
	}

	time.Sleep(time.Second * 35)

	var exist string
	go func() {
		var src string
		if err := chromedp.Run(ctx, chromedp.Text(`body > table > tbody > tr:nth-child(2) > td > div > div > div > div:nth-child(3) > div.row.cardBooks > div.col-sm-3.details-book-cover-container > div > a > span`,
			&src, chromedp.NodeVisible)); err != nil {
			log.Println(err)
			return
		}

		time.Sleep(3 * time.Second)
		exist = src
	}()

	time.Sleep(time.Second * 5)
	if exist == "已下载" {
		cancel()
		return errors.New(fmt.Sprintf("书籍已存在: url: %s", url))
	}

	fileDone := make(chan string, 1)

	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			//completed := "(unknown)"
			//if ev.TotalBytes != 0 {
			//	completed = fmt.Sprintf("%0.2f%%", ev.ReceivedBytes/ev.TotalBytes*100.0)
			//}
			//log.Printf("state: %s, completed: %s\n", ev.State.String(), completed)
			if ev.State == browser.DownloadProgressStateCompleted {
				fileDone <- "完成下载"
				close(fileDone)
			} else if ev.State == browser.DownloadProgressStateCanceled {
				fileDone <- fmt.Sprintf("下载失败: %s", url)
				close(fileDone)
			}
		}
	})

	if err := chromedp.Run(ctx,
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
			WithDownloadPath(dir).
			WithEventsEnabled(true),
		chromedp.Click(`body > table > tbody > tr:nth-child(2) > td > div > div > div > div:nth-child(3) > div:nth-child(2) > div.details-buttons-container.pull-left > div:nth-child(1) > div > a`, chromedp.NodeVisible)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("下载书籍失败 url:%s", url))

	}

	msg := <-fileDone
	time.Sleep(time.Second * 5)

	fmt.Println(msg)
	return nil
}
