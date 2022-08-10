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

type DangDangInfo struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string     `json:"title"`
	Intro     string     `json:"intro"`
	Author    string     `json:"author"`
	Publisher string     `json:"publisher"`
	Time      string     `json:"time"`
	ISBN      string     `json:"isbn"`
	Recommend string     `json:"recommend"`
	Desc      string     `json:"desc"`
	CreatedAt *time.Time `json:"createdAt"`
}

func main() {
	store.InitDB()
	store.InitRDS()

	store.DB.AutoMigrate(&DangDangInfo{})
	log := utils.NewLog("error.log")

	fmt.Println("start")
	for {
		l := store.RDS.LLen(context.Background(), "ddlink").Val()
		if l <= 0 {
			break
		}

		value := store.RDS.LPop(context.Background(), "ddlink")
		if value.Err() != nil {
			log.Error(errors.Wrap(value.Err(), `读取redis列表失败`))
			return
		}

		url := value.Val()

		func(url string) {

			ctx, cancel := chromedp.NewContext(
				context.Background(),
				chromedp.WithLogf(log.Printf),
			)
			defer cancel()

			ctx, cancel = context.WithTimeout(ctx, 40*time.Second)
			defer cancel()

			if err := chromedp.Run(ctx,
				utils.Setcookies(".dangdang.com", "sessionID", utils.DDSession, "secret_key", utils.DDSecret),
				chromedp.Navigate(url)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`打开url失败: %s`, url)))
				return
			}

			time.Sleep(20 * time.Second)

			var book DangDangInfo
			//标题
			if err := chromedp.Run(ctx,
				chromedp.AttributeValue(`#product_info > div.name_info > h1`, "title", &book.Title, nil)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取title失败: %s`, url)))
			}

			//介绍
			if err := chromedp.Run(ctx, chromedp.Text(`#author`, &book.Author, chromedp.NodeVisible)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取介绍失败: %s`, url)))
			}

			//作者
			if err := chromedp.Run(ctx, chromedp.Text(`#product_info > div.name_info > h2 > span.head_title_name`, &book.Intro, chromedp.NodeVisible)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取介绍失败: %s`, url)))
			}

			//出版社
			if err := chromedp.Run(ctx, chromedp.Text(`#product_info > div.messbox_info > span:nth-child(2) > a`, &book.Publisher, chromedp.NodeVisible)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取出版社失败: %s`, url)))
			}

			//时间
			if err := chromedp.Run(ctx, chromedp.Text(`#product_info > div.messbox_info > span:nth-child(3)`, &book.Time, chromedp.NodeVisible)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取出版社失败: %s`, url)))
			}

			//ISBN
			if err := chromedp.Run(ctx, chromedp.Text(`#detail_describe > ul > li:nth-child(5)`, &book.ISBN, chromedp.NodeVisible)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取出版社失败: %s`, url)))
			}

			//推荐
			if err := chromedp.Run(ctx, chromedp.Text(`#abstract > div.descrip`, &book.Recommend, chromedp.NodeVisible)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取推荐失败: %s`, url)))
			}

			//详情
			if err := chromedp.Run(ctx, chromedp.Text(`#content > div.descrip`, &book.Desc, chromedp.NodeVisible)); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`获取出版社失败: %s`, url)))
			}

			time.Sleep(3 * time.Second)
			cancel()

			if err := store.DB.Create(&book).Error; err != nil {
				log.Error(errors.Wrap(err, "创建记录失败"))
			}

			time.Sleep(1 * time.Minute)
			log.Info("\n")
		}(url)
	}

	fmt.Println("info end")
}
