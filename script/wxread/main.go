package main

import (
	"bufio"
	"context"
	"crawler/store"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type WxReadInfo struct {
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	SearchIndex int    `json:"search_idx"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	Cover       string `json:"cover"`
	ISBN        string `json:"isbn"`
	Publisher   string `json:"publisher"`
	Time        string `json:"publishTime"`
	Intro       string `json:"intro"`
	Desc        string `json:"desc"`
}

func main() {
	readLog()
	//store.InitDB()
	//log := utils.NewLog("error.log")
	//
	//f, err := os.Open("book.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer f.Close()
	//
	//r := bufio.NewReader(f)
	//
	//for {
	//	s, err := r.ReadString('\n')
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		} else {
	//			log.Error(err)
	//		}
	//		continue
	//	}
	//
	//	var data WxReadInfo
	//	if err := json.Unmarshal([]byte(s), &data); err != nil {
	//		log.Error(errors.Wrap(err, "解析失败"))
	//	}
	//	if err := store.DB.Create(&data).Error; err != nil {
	//		log.Error(errors.Wrap(err, fmt.Sprintf("创建记录失败 idx: %d", data.SearchIndex)))
	//	}
	//}

}

type logData struct {
	Msg string `json:"msg"`
}

func readLog() {
	store.InitRDS()

	f, err := os.Open("error.log")
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	r := bufio.NewReader(f)

	for {
		s, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
			}
			continue
		}

		var v logData
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			fmt.Println(err)
		}

		split := strings.Split(v.Msg, ":")
		url := "https:" + split[3]

		if err := store.RDS.RPush(context.Background(), "library-search", url).Err(); err != nil {
			fmt.Println(err)
		}
	}
}
