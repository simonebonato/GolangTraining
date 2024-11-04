package urlshort

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type UrlShort struct {
	gorm.Model
	URL string `gorm:"uniqueIndex:url"`
	Path string
}


var ShortenerDb, _ = gorm.Open(sqlite.Open("shortener.db"), &gorm.Config{})

func AddValues (db *gorm.DB) {
	// apply the migration automatically if needed
	db.AutoMigrate(&UrlShort{})

	// values 
	val_map := make(map[string]string)

	val_map["/test-db"] = "https://www.youtube.com/watch?v=mGtVzQ_d9oQ"
	val_map["/test-db2"] = "https://gorm.io/"

	for k, v := range(val_map) {
		targetUrlShort := UrlShort{URL: k}
		if res := db.First(&targetUrlShort, "url = ?", k); res.Error == gorm.ErrRecordNotFound {
			targetUrlShort.Path = v
			db.Create(&targetUrlShort)
		}
	}
}

func GetValues (db *gorm.DB, targetURL string) (path string, err bool) {
	// check if the value exists in the db
	var targetURLShort UrlShort
	if res := db.First(&targetURLShort, "url = ?", targetURL); res.Error != gorm.ErrRecordNotFound {
		return targetURLShort.Path, false
	}
	return "", true
}