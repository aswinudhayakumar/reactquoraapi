package candidates

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gbrlsnchs/jwt"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type Profiledetails struct {
	Posts     int64
	Questions int64
	Answers   int64
	Post      []Post
	Question  []Question
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

type CustomPayload struct {
	jwt.Payload
	Name     string
	Email    string
	Userid   uint
	Category string
	Profile  string
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

func Signup(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	fmt.Println("into signup")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 10")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var candidate Candidate
	var c Candidate
	err := vars.Decode(&candidate)

	fmt.Println(candidate)

	if err == nil {

		db.Where("Email = ?", candidate.Email).Find(&c)
		if c.ID > 0 {
			fmt.Fprintf(w, "fail")
		} else {
			db.Create(&Candidate{Email: candidate.Email, Password: candidate.Password, Googleid: candidate.Googleid, Name: candidate.Name, Profilepic: candidate.Profilepic, Category: candidate.Category})
			fmt.Fprintf(w, "success")
		}
	}

	defer db.Close()

}

func Signinwithgoogle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 11")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var candidate Candidate
	err := vars.Decode(&candidate)

	var c Candidate

	fmt.Println(candidate.Googleid, candidate.Email)

	if err == nil {
		db.Where("Email = ?", candidate.Email).Find(&c)

		if c.ID != 0 && c.Googleid == "" && candidate.Googleid != "" {
			c.Googleid = candidate.Googleid
			db.Save(&c)
			defer db.Close()
		}

		if c.ID != 0 && c.Googleid == candidate.Googleid {
			var now = time.Now()
			var hs256 = jwt.NewHS256([]byte("secret"))

			p := CustomPayload{
				Payload: jwt.Payload{
					Issuer:         "gbrlsnchs",
					Subject:        "someone",
					Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
					ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * 12 * time.Hour)),
					NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
					IssuedAt:       jwt.NumericDate(now),
					JWTID:          "foobar",
				},
				Name:     c.Name,
				Email:    c.Email,
				Userid:   c.ID,
				Category: c.Category,
				Profile:  c.Profilepic,
			}

			token, err := jwt.Sign(p, hs256)
			if err == nil {
				fmt.Fprintf(w, string(token))
				log.Printf("token = %s", token)
				fmt.Println(token)
			}
		} else {
			fmt.Fprintf(w, "null")
		}
	}

	defer db.Close()

}

func Signin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 12")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var candidate Candidate
	err := vars.Decode(&candidate)

	var c Candidate

	if err == nil {
		db.Where("Email = ?", candidate.Email).Find(&c)
		defer db.Close()
		if c.ID != 0 && c.Password == candidate.Password {
			var now = time.Now()
			var hs256 = jwt.NewHS256([]byte("secret"))

			p := CustomPayload{
				Payload: jwt.Payload{
					Issuer:         "gbrlsnchs",
					Subject:        "someone",
					Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
					ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * 12 * time.Hour)),
					NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
					IssuedAt:       jwt.NumericDate(now),
					JWTID:          "foobar",
				},
				Name:     c.Name,
				Email:    c.Email,
				Userid:   c.ID,
				Category: c.Category,
				Profile:  c.Profilepic,
			}

			token, err := jwt.Sign(p, hs256)
			if err != nil {
				// Handle error.
			}
			fmt.Fprintf(w, string(token))
			log.Printf("token = %s", token)
		} else {
			fmt.Fprintf(w, "null")
		}
	}

	defer db.Close()

}

func Verifypass(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 12")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var input Candidate
	err := vars.Decode(&input)
	var can Candidate

	if err == nil {
		db.Where("ID = ?", input.ID).Find(&can)
		if can.ID > 0 && can.Password == input.Password {
			fmt.Fprintf(w, "true")
		} else {
			fmt.Fprintf(w, "false")
		}
		fmt.Println(can.Password)
	}
	defer db.Close()
}

func Setpass(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 12")
	} else {
		fmt.Println("Connection Successfull")
	}

	defer db.Close()

	vars := json.NewDecoder(r.Body)
	var input Candidate
	err := vars.Decode(&input)
	var can Candidate

	if err == nil {
		db.Where("ID = ?", input.ID).Find(&can)
		can.Password = input.Password
		db.Save(can)
		fmt.Fprintf(w, "true")
	}

	fmt.Fprintf(w, "false")

}

func Getprofiledetails(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 12")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var input Candidate
	err := vars.Decode(&input)
	var output Profiledetails
	var post []Post
	var ques []Question
	var ans []Answer

	if err == nil {
		db.Where("Userid = ?", input.ID).Find(&post)
		db.Where("Userid = ?", input.ID).Find(&ques)
		db.Where("Userid = ?", input.ID).Find(&ans)

		output.Posts = int64(len(post))
		output.Questions = int64(len(ques))
		output.Answers = int64(len(ans))
		output.Post = post
		output.Question = ques

	}

	json.NewEncoder(w).Encode(output)

	defer db.Close()

}
