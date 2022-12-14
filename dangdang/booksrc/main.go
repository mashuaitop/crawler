package main

import (
	"context"
	"crawler/store"
	"crawler/utils"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"time"
)

func main() {
	store.InitRDS()

	log := utils.NewLog("error.log")

	fmt.Println("start")
	for {
		l := store.RDS.LLen(context.Background(), utils.RDSDBookNamekey).Val()
		if l <= 0 {
			break
		}

		value := store.RDS.LPop(context.Background(), utils.RDSDBookNamekey)
		if value.Err() != nil {
			log.Error(errors.Wrap(value.Err(), `读取redis列表失败`))
			return
		}

		name := value.Val()

		func(name string) {
			//opts := append(chromedp.DefaultExecAllocatorOptions[:],
			//	chromedp.Flag("headless", false),
			//)
			//
			//allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
			//defer cancel()

			ctx, cancel := chromedp.NewContext(
				context.Background(),
				chromedp.WithLogf(log.Printf),
			)
			defer cancel()

			ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			if err := chromedp.Run(ctx,
				chromedp.Navigate(fmt.Sprintf("http://search.dangdang.com/?key=%s&act=input", name))); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`搜索列表打开失败: %s`, name)))
				time.Sleep(5 * time.Second)
				return
			}

			time.Sleep(5 * time.Second)

			var href string
			if err := chromedp.Run(ctx, chromedp.AttributeValue(`#search_nature_rg ul li:nth-child(1) > a`, "href", &href, nil)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取链接失败: %s`, name)))
			}

			cancel()

			fmt.Printf("%s: %s \n", name, href)
			if href != "" {
				href = "http:" + href
				if err := store.RDS.RPush(context.Background(), "ddlink", href).Err(); err != nil {
					log.Error(errors.Wrap(err, fmt.Sprintf(`写入redis错误: %s`, href)))
				}
			}

			time.Sleep(time.Second * 20)
		}(name)
	}

	fmt.Println("finish")
}
