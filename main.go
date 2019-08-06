package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aswinudhayakumar/quoraapi/candidates"
	"github.com/aswinudhayakumar/quoraapi/post"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

var db *gorm.DB
var err error

type Candidate struct {
	gorm.Model

	Email      string
	Googleid   string
	Password   string
	Name       string
	Category   string
	Profilepic string
}

type Post struct {
	gorm.Model

	Category    string
	Title       string
	Name        string
	Image       string
	Postimage   string
	Description string
	Likes       int64
	Dislikes    int64
	Comments    int64
	Userid      int64
	Date        string
	Time        string
}

type Actions struct {
	gorm.Model

	Postid  int64
	Userid  int64
	Name    string
	Profile string
	Like    int64 `gorm:"default:0"`
	Dislike int64 `gorm:"default:0"`
	Comment string
}

type Question struct {
	gorm.Model

	Question  string
	Link      string
	Userimage string
	Name      string
	Userid    int64
	Category  string
	Answers   int64 `gorm:"default:0"`
}

type Answer struct {
	gorm.Model

	Userid     int64
	Questionid int64
	Profile    string
	Answer     string
	Name       string
	profile    string
}

type Notifications struct {
	gorm.Model

	Userid     int64
	Postuserid int64
	Postid     int64
	Image      string
	Name       string
	Message    string
}

func InitialMigration() {

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect")
	} else {
		fmt.Println("Migration successful")
	}

	defer db.Close()

	db.AutoMigrate(&Candidate{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Actions{})
	db.AutoMigrate(&Question{})
	db.AutoMigrate(&Answer{})
	db.AutoMigrate(&Notifications{})

}

func handlerequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/signup", candidates.Signup).Methods("POST")
	myRouter.HandleFunc("/signinwithgoogle", candidates.Signinwithgoogle).Methods("POST")
	myRouter.HandleFunc("/signin", candidates.Signin).Methods("POST")

	myRouter.HandleFunc("/like", post.Like).Methods("POST")
	myRouter.HandleFunc("/getliked", post.Getliked).Methods("POST")
	myRouter.HandleFunc("/getdisliked", post.Getdisliked).Methods("POST")
	myRouter.HandleFunc("/dislike", post.Dislike).Methods("POST")
	myRouter.HandleFunc("/comment", post.Comment).Methods("POST")

	myRouter.HandleFunc("/newpost", post.Newpost).Methods("POST")
	myRouter.HandleFunc("/uploadimage", post.Uploadimage).Methods("POST")
	myRouter.HandleFunc("/feed", post.Feed).Methods("POST")
	myRouter.HandleFunc("/getcomments", post.Getcomments).Methods("POST")

	myRouter.HandleFunc("/newquestion", post.Newquestion).Methods("POST")
	myRouter.HandleFunc("/getquestion", post.Getquestions).Methods("POST")
	myRouter.HandleFunc("/getsinglequestion", post.Getquestion).Methods("POST")
	myRouter.HandleFunc("/addcategory", post.Getanswers).Methods("POST")

	myRouter.HandleFunc("/addanswer", post.Addanswer).Methods("POST")
	myRouter.HandleFunc("/getanswers", post.Getanswers).Methods("POST")

	myRouter.HandleFunc("/setnotification", post.Setnotification).Methods("POST")
	myRouter.HandleFunc("/resetlikenotification", post.Resetlikenotification).Methods("POST")
	myRouter.HandleFunc("/resetcommentnotification", post.Resetcommentnotification).Methods("POST")
	myRouter.HandleFunc("/closenotofication", post.Closenotification).Methods("POST")
	myRouter.HandleFunc("/getnotification", post.Getnotification).Methods("POST")

	myRouter.HandleFunc("/verifypass", candidates.Verifypass).Methods("POST")
	myRouter.HandleFunc("/setpass", candidates.Setpass).Methods("POST")
	myRouter.HandleFunc("/getprofiledetails", candidates.Getprofiledetails).Methods("POST")

	myRouter.PathPrefix("/temp-images/").Handler(http.StripPrefix("/temp-images/", http.FileServer(http.Dir("temp-images"))))

	log.Fatal(http.ListenAndServe(":8123", cors.Default().Handler(myRouter)))
}

func main() {
	InitialMigration()
	handlerequests()
}
