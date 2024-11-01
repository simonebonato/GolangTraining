// here we are making use of the SQLite database library at "go get github.com/glebarez/sqlite"
// and the GORM library, to handle databases, at "go get -u gorm.io/gorm"
// tutorial from this great video here "https://www.youtube.com/watch?v=mGtVzQ_d9oQ"
package main

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// elements in the database, missing ID, createdAT, etc...
// it is by default in the gorm.Model
type Post struct {
	gorm.Model
	Title string
	Slug string `gorm:"uniqueIndex:idx_slug"` // we want to create this index for the slug
	Likes uint
}

func (p Post) String() string {
	return fmt.Sprintf("Post Title: %s, Slug: %s", p.Title, p.Slug)
}

var db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

func main() {
	// gorm will also take care of the migrations
	// first time we run "go run main.go" it will create the database 
	// and make the migrations
	// to see what is inside "sqlite3", then ".schema posts"
	db.AutoMigrate(&Post{})

	// else to make migrations manually do this 
	// in our case we added the "Likes" column after 
	// GORM is smart enough to see that there is a new col
	// and add it with the name "likes"
	// db.Migrator().AddColumn(&Post{}, "likes")

	// create the post and add it to the database
	// freshPost := createPost("Simo is the best", "simo-slug")
	// fmt.Println(freshPost)

	// to get a post from the database by the slug
	oldPost := getPost("simo-slug")
	fmt.Println(oldPost)

}

func createPost(title string, slug string) Post {
	
	// to add a value to the database, we first create the object in the memory
	// then pass the reference to the new post we created to the db 
	newPost := Post{Title: title, Slug: slug}
	if res := db.Create(&newPost); res.Error!=nil {
		panic(res.Error)
	}
	return newPost
}

func getPost(slug string) Post{
	targetPost := Post{Slug: slug}
	if res := db.First(&targetPost); res.Error != nil {
		panic(res.Error)
	}
	return targetPost
}