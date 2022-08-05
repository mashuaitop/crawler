package main

import (
	"context"
	"crawler/store"
	"crawler/utils"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
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
		if err = store.RDS.RPush(context.Background(), utils.RDSBookNamekey).Err(); err != nil {
			log.Error(errors.Wrap(err, fmt.Sprintf("写入书名失败: %s", name)))
		}
	}

	fmt.Println("finish ")
}
