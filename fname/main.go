package main

import (
	"crawler/utils"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	//store.InitRDS()

	log := utils.NewLog("error.log")

	dir, err := ioutil.ReadDir("/Users/mashuai/Downloads/book")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("./name.txt", os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, book := range dir {
		name := book.Name()
		if name[0] == '.' {
			continue
		}

		split := strings.Split(name, ".")
		bookName := split[0]
		fmt.Println(bookName)
		//if err = store.RDS.RPush(context.Background(), utils.RDSDBookNamekey, bookName).Err(); err != nil {
		//	log.Error(errors.Wrap(err, fmt.Sprintf("写入当当书名失败: %s", bookName)))
		//}
		_, err := io.WriteString(file, bookName+"\n")
		if err != nil {
			log.Error(err)
		}

	}

	fmt.Println("finish ")
}
