package methods

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
)

func SearchDetailHref(name string) (string, error)  {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36 Aoyou/cXRsNCdsM3s-T1c8SHhARZqOZNMwOHWB7sPpE_x2ULIWqtc__h71MI7ASQ==`),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Minute)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("https://zh.u1lib.org/s/%s?extensions[]=epub", name))); err != nil {
		return "",errors.Wrap(err, fmt.Sprintf(`搜索列表打开失败: %s`, name))
	}

	time.Sleep(time.Second * 20)
	var href string
	if err := chromedp.Run(ctx, chromedp.AttributeValue(`#searchResultBox > div:nth-child(2) > div > table > tbody > tr > td.itemCover > div > div > a`, "href", &href, nil)); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf(`获取书籍详情href失败: %s`, name))
	}

	time.Sleep(time.Second * 2)

	domain := "https://zh.u1lib.org"

	url := domain + href

	return url, nil
}

func SearchDetailImg(ctx context.Context) (string, error)  {
	var src string
	if err := chromedp.Run(ctx, chromedp.AttributeValue(`body > table > tbody > tr:nth-child(2) > td > div > div > div > div:nth-child(3) > div.row.cardBooks > div.col-sm-3.details-book-cover-container > div > a > div > img`, "src", &src, nil)); err != nil {
		return "", errors.Wrap(err, "获取书籍封面失败")
	}

	src = strings.Replace(src, "299", "", 1)

	return src, nil
}
