package server

import (
	"net/http"
	"log"
	"encoding/json"
	"time"
	"strconv"
	"encoding/hex"
	"crypto/rand"
	"calculator/resolver"
)

const (
	CORRECTLOGIN = "admin"
	CORRECTPASS  = "admin"
)

type Token struct {
	Token string `json:"token"`
	ExpirationDate int64 `json:"-"`
}

type UserReq struct {
	Login string `json:"login"`
	Passowrd string `json:"pass"`
}

type ComputeReq struct {
	Token string `json:"token"`
	Expression string `json:"expression"`
}

var USERTOKEN = make(map[string]*Token)
var TOKENS = make(map[string]*Token)

func generateToken() *Token {
	bytes := make([]byte,16)
	rand.Read(bytes)
	stringBuf := make([]byte,36)

	hex.Encode(stringBuf[0:8], bytes[0:4])
	stringBuf[8] = '-'
	hex.Encode(stringBuf[9:13], bytes[4:6])
	stringBuf[13] = '-'
	hex.Encode(stringBuf[14:18], bytes[6:8])
	stringBuf[18] = '-'
	hex.Encode(stringBuf[19:23], bytes[8:10])
	stringBuf[23] = '-'
	hex.Encode(stringBuf[24:], bytes[10:])


	return &Token{
		string(stringBuf),
		time.Now().Add(time.Minute * 15).Unix(),
	}
}

func Start() {
	http.HandleFunc("/login", authorization)
	http.HandleFunc("/compute",compute)

	err := http.ListenAndServe(":8080", nil)
	if err == nil {
		log.Fatal(err)
	} else {
		log.Printf("Server started successfully on port 8080")
	}

}

func authorization(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user UserReq
	err := decoder.Decode(&user)
	if err != nil {
		w.Write([]byte(`{"fault":"invalid request"}`))
		log.Printf("Bad request received: %v", err)
		return
	}
	log.Printf("Authentification request: %v\n",user)

	if user.Login == CORRECTLOGIN && user.Passowrd == CORRECTPASS {
		token, present := USERTOKEN[user.Login]

		if !present {
			token = generateToken()
			USERTOKEN[user.Login] = token
			TOKENS[token.Token] = token
			log.Printf("Token generated for %v : %v",user,token)
		} else {
			token.ExpirationDate = time.Now().Add(time.Minute * 15).Unix()
			log.Printf("Expiration date modified for user %v: %v",user,token)
		}

		byteRes, err := json.Marshal(token)
		if err == nil {
			w.Write(byteRes)
		} else {
			w.Write([]byte(`{"fault":"server side error"}`))
			log.Printf("Error while preparing the JSON answer: %v", err)
			return
		}
	} else {
		w.Write([]byte(`{"fault":"invalid login and password"}`))
		log.Printf("Wrong user/password given: %v", user)
	}
}

func compute(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var req ComputeReq
	err := decoder.Decode(&req)
	if err != nil {
		w.Write([]byte(`{"fault":"invalid request"}`))
		log.Printf("Wrong computation request: %v", err)
		return
	}
	log.Printf("Computation request: %v\n",req)

	token, present := TOKENS[req.Token]

	if !present || token.ExpirationDate < time.Now().Unix() {
		w.Write([]byte(`{"fault":"token invalid"}`))
		log.Printf("Request with invalid token: %v",present)
		return
	}

	val, err := resolver.Resolve(req.Expression)

	if err != nil {
		w.Write([]byte(`{"fault":"expression invalid"}`))
		log.Printf("Request with invalid expression or error during resolution: %v",err)
	} else {
		res := `{"result":` + val.String() + `}`
		log.Printf("Computation result: %v",res)
		w.Write([]byte(res))
	}
}

func (u UserReq) String() string {
	return "User(" + u.Login + ", " + u.Passowrd + ")"
}

func (t Token) String() string {
	return "Token(" + t.Token + ", " + strconv.FormatInt(t.ExpirationDate,10) + ")"
}

func (r ComputeReq) String() string {
	return "Request(" + r.Token + ", " + r.Expression + ")"
}