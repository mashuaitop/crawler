package methods

import "gorm.io/gorm"

type WxReadInfo struct {
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	SearchIndex int    `json:"searchIndex"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Cover       string `json:"cover"`
	ISBN        string `json:"isbn"`
	Publisher   string `json:"publisher"`
	Time        string `json:"time"`
	Category    string `json:"category"`
	Intro       string `json:"intro"`
	Desc        string `json:"desc"`
	IntroSync   bool   `json:"introSync" gorm:"default:false"`
	BookExist   bool   `json:"bookExist" gorm:"default:false"`
}

func WxBookName(db *gorm.DB) []string {
	var data []WxReadInfo
	db.Table(`wx_read_info`).Select(`title`).Where(`id > 911 and book_exist = false`).Order(`id`).Limit(100).Scan(&data)

	var names []string
	for _, row := range data {
		names = append(names, row.Title)
	}

	return names
}
