package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	dirPath := "/Users/mashuai/Downloads/bookcha/"

	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, book := range dir {
		name := book.Name()
		fmt.Println(name)
		if name[0] == '.' {
			continue
		}
		split := strings.Split(name, "【")
		if len(split) > 1 {
			newPath := dirPath + strings.TrimSpace(split[0]) + ".epub"
			if err = os.Rename(dirPath+name, newPath); err != nil {
				log.Println(err)
			}
		}
	}
}
