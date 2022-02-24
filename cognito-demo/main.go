package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/signup", signUp)

	port := os.Getenv("PORT")
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

// encoding 된 request body를 decode하기 위함
// decoding이란 JSON 문자열을 GO가 인식할 수 있는 값으로 바꾸는 것
type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func signUp(w http.ResponseWriter, r *http.Request) {
	// POST 요청이므로 request body가 함께 담겨 올 것이다.
	// parse the request body
	var req SignUpRequest
	// Decode 함수는 json값을 go value로 전환한 수 해당 변수에 담아준다.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// build a signup request

	// make the signup request

	w.Write([]byte("signup"))
}
