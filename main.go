package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aswinudhayakumar/quoraapi/candidates"

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

	Category   string
	Title      string
	Name       string
	Image      string
	Postimage  string
	Decription string
	Likes      int64
	Dislikes   int64
	Comments   int64
	Userid     int64
	Date       string
	Time       string
}

type Comments struct {
	gorm.Model

	Postid  int64
	Userid  int64
	Comment string
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
	db.AutoMigrate(&Comments{})

}

func handlerequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/signup", candidates.Signup).Methods("POST")
	myRouter.HandleFunc("/signinwithgoogle", candidates.Signinwithgoogle).Methods("POST")
	myRouter.HandleFunc("/signin", candidates.Signin).Methods("POST")

	myRouter.PathPrefix("/temp-images/").Handler(http.StripPrefix("/temp-images/", http.FileServer(http.Dir("temp-images"))))

	log.Fatal(http.ListenAndServe(":8123", cors.Default().Handler(myRouter)))
}

func main() {
	InitialMigration()
	handlerequests()
}
