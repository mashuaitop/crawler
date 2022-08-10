package main

import (
	"context"
	"crawler/library/model"
	"crawler/store"
	"crawler/utils"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	store.InitRDS()

	log := utils.NewLog("error.log")

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

			ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
			defer cancel()

			if err := chromedp.Run(ctx,
				chromedp.Navigate(fmt.Sprintf("https://zh.u1lib.org/s/%s?extensions[]=epub", name))); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`搜索列表打开失败: %s`, name)))
				time.Sleep(5 * time.Second)
				return
			}

			time.Sleep(time.Minute)

			var src string
			if err := chromedp.Run(ctx, chromedp.AttributeValue(`#searchResultBox > div:nth-child(2) > div > table > tbody > tr > td.itemCover > div > div > a > img`, "src", &src, nil)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取图片src 失败: %s`, name)))
				return
			}

			cancel()

			src = strings.Replace(src, "100", "", 1)
			fmt.Printf("%s: %s \n", name, src)
			if src != "" {
				var info model.ImgInfo
				info.Name = name
				info.Src = src

				data, err := json.Marshal(&info)
				if err != nil {
					log.Error(err)
					return
				}

				if err = store.RDS.RPush(context.Background(), utils.RDSBookSrckey, string(data)).Err(); err != nil {
					log.Error(errors.Wrap(err, fmt.Sprintf("写入图书馆书籍图片失败: %s", name)))
				}
			}

		}(name)
	}

	fmt.Println("finish")
}
