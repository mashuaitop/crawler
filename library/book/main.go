package main

import (
	"context"
	"crawler/library/methods"
	"crawler/store"
	"crawler/utils"
	"fmt"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"sync"
	"time"
)

func main() {
	store.InitDB()
	store.InitRDS()

	//searchBook()

	downloadBook()
}

func searchBook() {
	names := methods.WxBookName(store.DB)

	searchChannel := make(chan struct{}, 8)
	writeCh := make(chan string, 12)
	var wg sync.WaitGroup

	go func() {
		for url := range writeCh {
			fmt.Printf("url: %s \n", url)
			if err := store.RDS.RPush(context.Background(), "library-search", url).Err(); err != nil {
				fmt.Println(err)
			}
		}
	}()

	for _, name := range names {
		searchChannel <- struct{}{}
		wg.Add(1)
		go func(name string) {
			fmt.Printf("bookName: %s \n", name)
			url, err := methods.SearchDetailHref(name)
			if err != nil {
				fmt.Println(err)
			}

			if url != "" {
				writeCh <- url
			}
			defer wg.Done()
			<-searchChannel

			return
		}(name)
	}

	wg.Wait()

	fmt.Println("ok")
}

func downloadBook() {
	userID := "24968293"
	userKey := "b09825b8a4a84888b43cc0d6b4820306"
	bookPath := "/Volumes/data/cbook/"
	log := utils.NewLog("error.log")

	//url := "https://zh.u1lib.org/book/13956997/40e05c"
	dlCh := make(chan struct{}, 8)
	var wg sync.WaitGroup

	for i := 0; i < 200; i++ {
		length := store.RDS.LLen(context.Background(), "library-search").Val()
		if length == 0 {
			fmt.Println("没有数据了")
			break
		}

		fmt.Println(i)
		dlCh <- struct{}{}
		wg.Add(1)

		url := store.RDS.LPop(context.Background(), "library-search").Val()
		if url == "" {
			continue
		}

		go func(url string) {
			defer wg.Done()
			if err := methods.DownloadBook(userID, userKey, url, bookPath); err != nil {
				log.Error(err)
			}

			<-dlCh
		}(url)
	}

	wg.Wait()

	fmt.Println("任务完成")
}

func old() {
	userID := "24783016"
	userKey := "224bfbc05c76f8d995b3f4ffa46ce4e3"
	imgPath := "/Users/mashuai/Downloads/bookimg/"
	bookPath := "/Users/mashuai/Downloads/book/"

	log := utils.NewLog("error.log")

	names, err := utils.ReadLine("./1.txt")
	if err != nil {
		log.Error(err)
		return
	}

	for _, bookName := range names {
		func() {
			url, err := methods.SearchDetailHref(bookName)
			if err != nil {
				log.Error(err)
				return
			}
			fmt.Println("url: ", url)

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
				utils.Setcookies(".u1lib.org", "remix_userid", userID, "remix_userkey", userKey),
				chromedp.Navigate(url)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf("打开书籍详情页失败: name: %s; url: %s", bookName, url)))
				return
			}

			time.Sleep(time.Second * 20)

			imgDone := make(chan string, 1)
			go func() {
				defer close(imgDone)
				src, err := methods.SearchDetailImg(ctx)
				if err != nil {
					log.Error(errors.Wrap(err, fmt.Sprintf("name: %s", bookName)))
					return
				}

				if err = methods.DownloadImg(bookName, src, imgPath); err != nil {
					log.Error(err)
				}
				imgDone <- "ok"
				return
			}()

			fileDone := make(chan string, 1)
			chromedp.ListenTarget(ctx, func(v interface{}) {
				if ev, ok := v.(*browser.EventDownloadProgress); ok {
					completed := "(unknown)"
					if ev.TotalBytes != 0 {
						completed = fmt.Sprintf("%0.2f%%", ev.ReceivedBytes/ev.TotalBytes*100.0)
					}
					log.Printf("state: %s, completed: %s\n", ev.State.String(), completed)
					if ev.State == browser.DownloadProgressStateCompleted {
						fileDone <- "ok"
						close(fileDone)
					}
				}
			})

			if err := chromedp.Run(ctx,
				browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
					WithDownloadPath(bookPath).
					WithEventsEnabled(true),
				chromedp.Click(`body > table > tbody > tr:nth-child(2) > td > div > div > div > div:nth-child(3) > div:nth-child(2) > div.details-buttons-container.pull-left > div:nth-child(1) > div > a`, chromedp.NodeVisible)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf("下载书籍失败 url: name: %s; %s", bookName, url)))
				return
			}

			<-imgDone
			<-fileDone
			time.Sleep(time.Second * 20)
			log.Info(fmt.Sprintf("%s 完成下载", bookName))
			return
		}()
	}

	fmt.Println("end")
}
