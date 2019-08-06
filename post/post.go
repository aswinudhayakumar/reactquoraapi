package post

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type Answer struct {
	gorm.Model

	Userid     int64
	Questionid int64
	Profile    string
	Answer     string
	Name       string
	profile    string
}

type Candidate struct {
	gorm.Model

	Email      string
	Googleid   string
	Password   string
	Name       string
	Category   string
	Profilepic string
}

type Cat struct {
	Category string
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

type Post struct {
	gorm.Model

	Category    string
	Title       string
	Name        string
	Image       string
	Postimage   string
	Description string
	Likes       int64 `gorm:"default:0"`
	Dislikes    int64 `gorm:"default:0"`
	Comments    int64 `gorm:"default:0"`
	Userid      int64
	Date        string
	Time        string
}

type Actions struct {
	gorm.Model

	Postid  int64
	Userid  int64
	Profile string
	Name    string
	Like    int64 `gorm:"default:0"`
	Dislike int64 `gorm:"default:0"`
	Comment string
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

func Like(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 1")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var bodyaction Actions
	err := vars.Decode(&bodyaction)
	var post Post
	var action Actions

	if err == nil {

		db.Where("Postid = ? AND Userid = ? ", bodyaction.Postid, bodyaction.Userid).Find(&action)
		db.Where("Id = ?", action.Postid).Find(&post)

		if action.ID > 0 {

			if action.Like == 1 {
				post.Likes = post.Likes - 1
				action.Like = 0
				db.Save(post)
			} else {
				post.Likes = post.Likes + 1
				if action.Dislike == 1 {
					action.Dislike = 0
					post.Dislikes = post.Dislikes - 1
				}
				action.Like = 1
				db.Save(post)
			}
			db.Save(action)
		} else {

			db.Create(&Actions{Name: bodyaction.Name, Like: 1, Profile: bodyaction.Profile, Userid: bodyaction.Userid, Postid: bodyaction.Postid, Dislike: 0})

			db.Where("Postid = ? AND Userid = ? ", bodyaction.Postid, bodyaction.Userid).Find(&action)
			db.Where("Id = ?", action.Postid).Find(&post)

			post.Likes = post.Likes + 1
			db.Save(post)
		}
	} else {
		fmt.Println(err)
	}

	defer db.Close()

}

func Dislike(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 2")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var bodyaction Actions
	err := vars.Decode(&bodyaction)
	var action Actions
	var post Post

	if err == nil {

		db.Where("Postid = ? AND Userid = ? ", bodyaction.Postid, bodyaction.Userid).Find(&action)
		db.Where("Id = ?", action.Postid).Find(&post)

		if action.ID > 0 {
			if action.Dislike == 1 {
				post.Dislikes = post.Dislikes - 1
				action.Dislike = 0
				db.Save(post)
			} else {
				post.Dislikes = post.Dislikes + 1
				if action.Like == 1 {
					action.Like = 0
					post.Likes = post.Likes - 1
				}
				action.Dislike = 1
				db.Save(post)
			}
			db.Save(action)
		} else {
			db.Create(&Actions{Name: bodyaction.Name, Like: 0, Profile: bodyaction.Profile, Userid: bodyaction.Userid, Postid: bodyaction.Postid, Dislike: 1})

			db.Where("Postid = ? AND Userid = ? ", bodyaction.Postid, bodyaction.Userid).Find(&action)
			db.Where("Id = ?", action.Postid).Find(&post)

			post.Dislikes = post.Dislikes + 1
			db.Save(post)
		}
	} else {
		fmt.Println(err)
	}

	defer db.Close()

}

func Comment(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 3")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var bodyaction Actions
	err := vars.Decode(&bodyaction)
	var post Post

	fmt.Println("bc ", bodyaction.Profile)

	if err == nil {
		var action Actions
		db.Where("Postid = ? AND Userid = ? ", bodyaction.Postid, bodyaction.Userid).Find(&action)
		db.Where("Id = ?", action.Postid).Find(&post)

		if action.ID > 0 {

			if bodyaction.Comment == "delpconf.conf" {
				action.Comment = ""
				post.Comments = post.Comments - 1
				db.Save(post)
				db.Save(action)
			} else {
				if action.Comment == "" {
					post.Comments = post.Comments + 1
					db.Save(post)
				}
				action.Comment = bodyaction.Comment
				db.Save(action)
			}
		} else {
			fmt.Println("this")
			db.Create(&Actions{Name: bodyaction.Name, Like: 0, Profile: bodyaction.Profile, Userid: bodyaction.Userid, Postid: bodyaction.Postid, Dislike: 0, Comment: bodyaction.Comment})

			db.Where("Postid = ? AND Userid = ? ", bodyaction.Postid, bodyaction.Userid).Find(&action)
			db.Where("Id = ?", action.Postid).Find(&post)

			post.Comments = post.Comments + 1
			db.Save(post)

		}
	} else {
		fmt.Println(err)
	}
	defer db.Close()
}

func Newpost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 4")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var newpost Post
	err := vars.Decode(&newpost)
	if err == nil {
		fmt.Println("Hello")
		db.Create(&Post{
			Category:    newpost.Category,
			Name:        newpost.Name,
			Title:       newpost.Title,
			Image:       newpost.Image,
			Postimage:   "",
			Description: newpost.Description,
			Likes:       0,
			Dislikes:    0,
			Comments:    0,
			Userid:      newpost.Userid,
			Date:        newpost.Date,
			Time:        newpost.Time,
		})
	} else {
		fmt.Println(err)
	}
	defer db.Close()
}

func Feed(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 5")
	} else {
		fmt.Println("Connection Successfull")
	}

	var cat Cat
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&cat)
	var post []Post
	defer db.Close()
	if err == nil {

		if cat.Category == "Feed" {
			db.Find(&post)
		} else if cat.Category == "Photography" {
			db.Where("Category = ?", cat.Category).Find(&post)
		} else if cat.Category == "Science" {
			db.Where("Category = ?", cat.Category).Find(&post)
		} else if cat.Category == "Literature" {
			db.Where("Category = ?", cat.Category).Find(&post)
		} else if cat.Category == "Health" {
			db.Where("Category = ?", cat.Category).Find(&post)
		} else if cat.Category == "Cooking" {
			db.Where("Category = ?", cat.Category).Find(&post)
		} else if cat.Category == "Music" {
			db.Where("Category = ?", cat.Category).Find(&post)
		} else if cat.Category == "Sports" {
			db.Where("Category = ?", cat.Category).Find(&post)
		} else if cat.Category == "Movies" {
			db.Where("Category = ?", cat.Category).Find(&post)
		}
	}

	json.NewEncoder(w).Encode(post)

}

func Getcomments(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 6")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Actions
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	var comment []Actions

	fmt.Println(input.Postid)

	if err == nil {
		db.Where("Postid = ? AND Comment != ''", input.Postid).Find(&comment)
	} else {
		fmt.Println(err)
	}

	fmt.Println(comment)

	json.NewEncoder(w).Encode(comment)
	defer db.Close()
}

func Uploadimage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	r.ParseMultipartForm(5 << 20)

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 7")
	} else {
		fmt.Println("Connection Successfull")
	}

	var post Post
	db.Last(&post)

	file, handler, err := r.FormFile("image")

	if err != nil {
		fmt.Println("Error uploading file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Println(handler.Filename)
	fmt.Println(handler.Size)
	fmt.Println(handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	ty := strings.Split(http.DetectContentType(fileBytes), "/")
	typ := ty[len(ty)-1:][0]
	fmt.Println(typ)
	tempfile, err := ioutil.TempFile("temp-images", "*."+typ)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(tempfile.Name())

	var link = tempfile.Name()

	tempfile.Write(fileBytes)

	post.Postimage = link
	db.Save(post)

	defer tempfile.Close()
	defer db.Close()

}

func Getliked(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 8")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Actions
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	var like Actions

	fmt.Println(input.Postid, input.Userid)

	if err == nil {
		db.Where("Postid = ? AND Userid = ?", input.Postid, input.Userid).Find(&like)
		fmt.Println(&like)
	}

	json.NewEncoder(w).Encode(like)
	defer db.Close()

}

func Getdisliked(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Actions
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	var like Actions

	if err == nil {
		db.Where("Postid = ? AND Userid = ?").Find(&like)
	}

	json.NewEncoder(w).Encode(like)
	defer db.Close()
}

func Newquestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Question
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	if err == nil {
		db.Create(&Question{Category: input.Category, Userid: input.Userid, Name: input.Name, Question: input.Question, Link: input.Link, Userimage: input.Userimage})
	}

	defer db.Close()
}

func Addcategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Candidate
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	if err == nil {
		var can Candidate
		db.Where("Email = ?", input.Email).Find(&can)
		can.Category = input.Category
		db.Save(can)
	}

	defer db.Close()

}
func Getquestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var ques []Question
	db.Find(&ques)

	json.NewEncoder(w).Encode(ques)

	defer db.Close()
}

func Getquestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Question
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)
	var ques Question

	if err == nil {
		db.Where("Id = ?", input.ID).Find(&ques)
	}

	json.NewEncoder(w).Encode(ques)

	defer db.Close()
}

func Addanswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Answer
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)
	var ques Question

	if err == nil {
		db.Create(&Answer{Questionid: input.Questionid, Profile: input.Profile, Name: input.Name, Answer: input.Answer, Userid: input.Userid, profile: input.profile})
		db.Where("Id = ?", input.Questionid).Find(&ques)
		ques.Answers = ques.Answers + 1
		db.Save(ques)
	}

	defer db.Close()

}

func Getanswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Question
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)
	var ans []Answer
	fmt.Println(input.ID)

	if err == nil {
		db.Where("Questionid = ?", input.ID).Find(&ans)
	}

	json.NewEncoder(w).Encode(ans)

	defer db.Close()
}
func Setnotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Notifications
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	if err == nil {
		db.Create(&Notifications{Userid: input.Userid, Image: input.Image, Postuserid: input.Postuserid, Postid: input.Postid, Message: input.Message, Name: input.Name})
		fmt.Println("Notification done")
	}

	defer db.Close()
}
func Resetlikenotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	fmt.Println("into reset liked")

	var input Notifications
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	if err == nil {
		db.Where("Postuserid = ? AND Postid = ? AND Message = 'like'", input.Postuserid, input.Postid).Unscoped().Delete(Notifications{})
		fmt.Println("Notification deleted")
	}

	defer db.Close()
}
func Resetcommentnotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Notifications
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	if err == nil {
		db.Where("Postuserid = ? AND Postid = ? AND Message = 'commented'", input.Postuserid, input.Postid).Unscoped().Delete(Notifications{})
	}

	defer db.Close()

}

func Closenotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Notifications
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)

	if err == nil {
		db.Where("ID = ?", input.ID).Unscoped().Delete(Notifications{})
	}
	defer db.Close()
}

func Getnotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 9")
	} else {
		fmt.Println("Connection Successfull")
	}

	var input Notifications
	vars := json.NewDecoder(r.Body)
	err := vars.Decode(&input)
	var notification []Notifications

	if err == nil {
		db.Where("Postuserid = ?", input.Postuserid).Find(&notification)
	}

	json.NewEncoder(w).Encode(notification)

	defer db.Close()
}
