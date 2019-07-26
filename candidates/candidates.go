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

type CustomPayload struct {
	jwt.Payload
	Name     string
	Email    string
	Userid   uint
	Category string
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

func Signup(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 1")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var candidate Candidate
	err := vars.Decode(&candidate)

	if err != nil {
		db.Create(&Candidate{Email: candidate.Email, Password: candidate.Password, Googleid: candidate.Googleid, Name: candidate.Name, Profilepic: candidate.Profilepic, Category: candidate.Category})
	}

	defer db.Close()

}

func Signinwithgoogle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 1")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var candidate Candidate
	err := vars.Decode(&candidate)

	var c Candidate

	if err != nil {
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

func Signin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=quora password=helloworld sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect 1")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := json.NewDecoder(r.Body)
	var candidate Candidate
	err := vars.Decode(&candidate)

	var c Candidate

	if err != nil {
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
