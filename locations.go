package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	clientSecret string
)

const clientID = "1YDQsQs35jh33XfAPL8T0KW5fz7jizOZ"

func main() {
	clientSecret = os.Getenv("MOVESSECRET")
	if clientSecret == "" {
		log.Fatal("$MOVESSECRET must be set")
	}

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8632"
	}

	err := migrate()
	if err != nil {
		log.Fatalf("migration failed, %v", err)
	}

	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/api/hello", handleHello).Methods("GET")
	r.HandleFunc("/api/authorize", handleAuthorize).Methods("GET")
	r.HandleFunc("/api/authcodeexchange", handleAuthCodeExchange).Methods("GET")
	r.HandleFunc("/api/users/{id}/counts", handleCounts).Methods("GET")
	r.HandleFunc("/api/users/{id}/sync", handleSyncGet).Methods("GET")
	r.HandleFunc("/api/users/{id}/sync", handleSync).Methods("POST")

	fmt.Println("Starting HTTP server")
	err = http.ListenAndServe(":"+port, r)
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
	fmt.Println("auth code", code)

	token, err := getAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	fmt.Println("access token", token)

	firstDate, err := getFirstDate(token.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	fmt.Println("first date", firstDate)

	err = createUser(token, firstDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func getAccessToken(code string) (accessToken, error) {
	url := fmt.Sprintf(
		"https://api.moves-app.com/oauth/v1/access_token?grant_type=authorization_code&code=%v&client_id=%v&client_secret=%v",
		code,
		clientID,
		clientSecret)

	resp, err := http.Post(url, "", nil)
	if err != nil {
		return accessToken{}, err
	}

	defer resp.Body.Close()
	token := accessToken{}
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return accessToken{}, err
	}

	return token, nil
}

func getFirstDate(accessToken string) (string, error) {
	url := fmt.Sprintf(
		"https://api.moves-app.com/api/1.1/user/profile?access_token=%v",
		accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	fmt.Println("profile status", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var dat map[string]interface{}

	err = json.Unmarshal(body, &dat)
	if err != nil {
		return "", err
	}

	profile := dat["profile"].(map[string]interface{})
	firstDate := profile["firstDate"].(string)

	return firstDate, nil
}

func handleSyncGet(w http.ResponseWriter, req *http.Request) {
	parmUserID := mux.Vars(req)["id"]
	if parmUserID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parmUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, uss, err := loadUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(uss)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to serialize response: %v", err), http.StatusInternalServerError)
		return
	}
}

func handleSync(w http.ResponseWriter, req *http.Request) {
	parmUserID := mux.Vars(req)["id"]
	if parmUserID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parmUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usr, uss, err := loadUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	i := 0
	for uss.SyncedThroughDate < time.Now().AddDate(0, 0, -7).Format("20060102") && i < 60 {
		i++

		ds, err := getDailyStoryline(usr.AccessToken, uss.SyncedThroughDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bcs := extractBreadcrumbs(ds)

		err = saveBreadcrumbs(usr.UserID, bcs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		syncDt, err := time.Parse("20060102", uss.SyncedThroughDate)
		syncDt = syncDt.AddDate(0, 0, 1)
		uss.SyncedThroughDate = syncDt.Format("20060102")

		err = saveUserSyncedThroughDate(usr.UserID, uss.SyncedThroughDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(uss)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to serialize response: %v", err), http.StatusInternalServerError)
		return
	}
}

func handleCounts(w http.ResponseWriter, req *http.Request) {
	parmUserID := mux.Vars(req)["id"]
	if parmUserID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parmUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parmStartDate := req.URL.Query().Get("startDate")
	if parmStartDate != "" {
		_, err = strconv.Atoi(parmStartDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	parmEndDate := req.URL.Query().Get("endDate")
	if parmEndDate != "" {
		_, err = strconv.Atoi(parmEndDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	ccs, err := loadCounts(userID, parmStartDate, parmEndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(ccs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to serialize response: %v", err), http.StatusInternalServerError)
		return
	}
}

func getDailyStoryline(bearerToken string, date string) (dailyStoryline, error) {
	url := fmt.Sprintf(
		"https://api.moves-app.com/api/1.1/user/storyline/daily/%v?trackPoints=true&access_token=%v",
		date,
		bearerToken)

	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return dailyStoryline{}, err
	}

	fmt.Println(resp.StatusCode)

	defer resp.Body.Close()
	storylines := []dailyStoryline{}
	err = json.NewDecoder(resp.Body).Decode(&storylines)
	if err != nil {
		return dailyStoryline{}, err
	}

	return storylines[0], nil
}

func extractBreadcrumbs(ds dailyStoryline) []breadcrumb {
	breadcrumbs := []breadcrumb{}

	for _, segment := range ds.Segments {
		for _, activity := range segment.Activities {
			for _, trackPoint := range activity.TrackPoints {
				bc := breadcrumb{
					Coordinate: coordinate{
						Lat: trackPoint.Lat,
						Lon: trackPoint.Lon,
					},
					Time: trackPoint.Time,
				}
				breadcrumbs = append(breadcrumbs, bc)
			}
		}
	}

	return breadcrumbs
}
