package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	dirPath := "/Users/mashuai/Downloads/book/"

	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, book := range dir {
		idx := strings.LastIndex(book.Name(), "crdownload")
		if idx > 0 {
			if err = os.Remove(dirPath + book.Name()); err != nil {
				fmt.Println(err)
			}
		}
	}

}
