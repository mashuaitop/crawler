package main

import (
	"context"
	"crawler/library/model"
	"crawler/store"
	"crawler/utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	store.InitRDS()

	path := "/Users/mashuai/Downloads/bookimg/"
	log := utils.NewLog("error.log")

	fmt.Println("start")
	for {
		l := store.RDS.LLen(context.Background(), utils.RDSBookSrckey).Val()
		if l <= 0 {
			break
		}

		value := store.RDS.LPop(context.Background(), utils.RDSBookSrckey)
		if value.Err() != nil {
			log.Error(errors.Wrap(value.Err(), `读取redis列表失败`))
			return
		}

		data := value.Val()

		var info model.ImgInfo
		if err := json.Unmarshal([]byte(data), &info); err != nil {
			log.Error(errors.Wrap(err, `json unmarshal err`))
			return
		}

		func(info *model.ImgInfo) {
			src := info.Src
			name := info.Name
			resp, err := http.Get(info.Src)
			if err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`http get err url: %s`, src)))
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`read body err url: %s`, src)))
				return
			}

			idx := strings.LastIndex(src, "/")
			var ext string
			spilt := strings.Split(src[idx:], ".")
			if len(spilt) > 0 {
				ext = spilt[1]
			}

			imgName := fmt.Sprintf(`%s%s.%s`, path, name, ext)
			fmt.Println(imgName)
			if err = ioutil.WriteFile(imgName, body, 0666); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf(`download  err url: %s`, src)))
				return
			}
		}(&info)
		break
	}

	fmt.Println("finish")
}
