package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	clientSecret string
)

const clientID = "1YDQsQs35jh33XfAPL8T0KW5fz7jizOZ"

type accessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

func main() {
	clientSecret = os.Getenv("MOVESSECRET")
	if clientSecret == "" {
		log.Fatal("$MOVESSECRET must be set")
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/hello", handleHello).
		Methods("GET")
	r.HandleFunc("/api/authorize", handleAuthorize).
		Methods("GET")
	r.HandleFunc("/api/authcodeexchange", handleAuthCodeExchange).
		Methods("GET")

	err := http.ListenAndServe(":8632", r)
	if err != nil {
		log.Fatalf("Cannot listen and serve, %v", err)
	} else {
		fmt.Printf("Done\n")
	}
}

func handleHello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello, world")
}

func handleAuthorize(w http.ResponseWriter, req *http.Request) {
	http.Redirect(
		w,
		req,
		fmt.Sprintf(
			"https://api.moves-app.com/oauth/v1/authorize?response_type=code&client_id=%v&scope=activity location",
			clientID),
		http.StatusFound)
}

func handleAuthCodeExchange(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code is a required querystring parameter", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, code)

	token, err := getAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getAccessToken(code string) (accessToken, error) {
	url := fmt.Sprintf(
		"https://api.moves-app.com/oauth/v1/access_token?grant_type=authorization_code&code=%v&client_id=%v&client_secret=%v",
		code,
		clientID,
		clientSecret)

	resp, err := http.Get(url)
	if err != nil {
		return accessToken{}, err
	}

	defer resp.Body.Close()
	token := accessToken{}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return accessToken{}, err
	}

	fmt.Println(token)

	return token, nil
}
