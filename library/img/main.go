package main

import (
	"context"
	"crawler/library/methods"
	"crawler/store"
	"crawler/utils"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"time"
)

func main()  {
	store.InitRDS()

	log := utils.NewLog("error.log")
	imgPath := "/Users/mashuai/Downloads/bookimg/"

	dir, err := ioutil.ReadDir("/Users/mashuai/Downloads/book")
	if err != nil {
		log.Fatal(err)
	}

	for _, book := range dir {
		name := book.Name()
		if name[0] == '.' {
			continue
		}

		split := strings.Split(name, ".")
		if len(split) > 0 {
			bookName := split[0]
			idx := strings.LastIndex(bookName, "z-lib")
			if idx != -1 {
				bookName = strings.TrimSpace(bookName[:idx])
			}

			fmt.Println(bookName)
			if err = store.RDS.RPush(context.Background(), utils.RDSZBookNamekey, bookName).Err(); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf("写入图书馆书名失败: %s", bookName)))
			}
		}
	}

	fmt.Println("start")

	for {
		l := store.RDS.LLen(context.Background(), utils.RDSZBookNamekey).Val()
		if l <= 0 {
			break
		}

		value := store.RDS.LPop(context.Background(), utils.RDSZBookNamekey)
		if value.Err() != nil {
			log.Error(errors.Wrap(value.Err(), `读取redis列表失败`))
			return
		}

		name := value.Val()

		func(name string) {
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

			ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
			defer cancel()

			if err := chromedp.Run(ctx,
				chromedp.Navigate(fmt.Sprintf("https://zh.u1lib.org/s/%s?extensions[]=epub", name))); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`搜索列表打开失败: %s`, name)))
				return
			}

			time.Sleep(time.Second * 30)

			var src string
			if err := chromedp.Run(ctx, chromedp.AttributeValue(`#searchResultBox > div:nth-child(2) > div > table > tbody > tr > td.itemCover > div > div > a > img`, "src", &src, nil)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取图片src 失败: %s`, name)))
				return
			}

			cancel()

			if src != "" {
				src = strings.Replace(src, "100", "", 1)
				if err = methods.DownloadImg(name, src, imgPath); err != nil {
					log.Error(err)
				}
			}
			time.Sleep(time.Second * 3)
		}(name)
	}

	fmt.Println("end")
}
