package main

import (
	"context"
	"crawler/utils"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"time"
)

type DouBanInfo struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string `json:"title"`
	Intro     string `json:"intro"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Time      string `json:"time"`
	ISBN      string `json:"isbn"`
	Desc      string `json:"desc"`
}

func main() {
	log := utils.NewLog("error.log")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36 Aoyou/cXRsNCdsM3s-T1c8SHhARZqOZNMwOHWB7sPpE_x2ULIWqtc__h71MI7ASQ==`),
		chromedp.Flag("headless", false),
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

	url := "https://book.douban.com/subject/35490038/"

	if err := chromedp.Run(ctx,
		utils.Setcookies(".douban.com",
			"gr_user_id", "d5bcc6d0-bc52-4b96-bf70-19b1e26f34d8",
			"bid", "BawKkkK5NpA",
		),
		chromedp.Navigate(url)); err != nil {
		log.Error(errors.Wrap(err, fmt.Sprintf(`打开url失败: %s`, url)))
		return
	}

	time.Sleep(30 * time.Second)

	var book DouBanInfo

	if err := chromedp.Run(ctx,
		chromedp.TextContent(`#info`, &book.Intro, chromedp.NodeVisible)); err != nil {
		log.Error(errors.Wrap(err, fmt.Sprintf(`获取介绍失败: %s`, url)))
	}

	//标题
	if err := chromedp.Run(ctx,
		chromedp.Text(`#wrapper > h1 > span`, &book.Title, chromedp.NodeVisible)); err != nil {
		log.Error(errors.Wrap(err, fmt.Sprintf(`获取title失败: %s`, url)))
	}

	//作者
	if err := chromedp.Run(ctx, chromedp.Text(`#info > span:nth-child(1) > a`, &book.Author, chromedp.NodeVisible)); err != nil {
		log.Error(errors.Wrap(err, fmt.Sprintf(`获取作者失败: %s`, url)))
	}

	//出版社
	if err := chromedp.Run(ctx, chromedp.Text(`#info > a:nth-child(4)`, &book.Publisher, chromedp.NodeVisible)); err != nil {
		log.Error(errors.Wrap(err, fmt.Sprintf(`获取出版社失败: %s`, url)))
	}

	//时间
	if err := chromedp.Run(ctx, chromedp.Text(`#info > span:nth-child(11)`, &book.Time, chromedp.NodeVisible)); err != nil {
		log.Error(errors.Wrap(err, fmt.Sprintf(`获取发布时间失败: %s`, url)))
	}

	//ISBN
	if err := chromedp.Run(ctx, chromedp.Text(`#info > span:nth-child(19)`, &book.ISBN, chromedp.NodeVisible)); err != nil {
		log.Error(errors.Wrap(err, fmt.Sprintf(`获取ISBN失败: %s`, url)))
	}

	//详情
	if err := chromedp.Run(ctx, chromedp.Text(`#link-report > span.all.hidden > div > div`, &book.Desc, chromedp.NodeVisible)); err != nil {
		log.Error(errors.Wrap(err, fmt.Sprintf(`获取详情失败: %s`, url)))
	}

	time.Sleep(10 * time.Second)
	cancel()

	fmt.Printf("%+v", book)
}
