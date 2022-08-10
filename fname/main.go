package main

import (
	"context"
	"crawler/store"
	"crawler/utils"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
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
		bookName := split[0]
		fmt.Println(bookName)
		if err = store.RDS.RPush(context.Background(), utils.RDSDBookNamekey, bookName).Err(); err != nil {
			log.Error(errors.Wrap(err, fmt.Sprintf("写入当当书名失败: %s", bookName)))
		}
	}

	fmt.Println("finish ")
}
