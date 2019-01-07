package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8632"
	}

	err := migrate()
	if err != nil {
		log.Fatalf("migration failed, %v", err)
	}

	r := mux.NewRouter()

	// This will serve files under http://<site>/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
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
	start := time.Now()

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
	for uss.SyncedThroughDate < time.Now().Format("20060102") &&
		i < 60 &&
		time.Now().Sub(start).Seconds() < 27 {

		i++

		ds, err := getData(usr.FollowMeeKey, usr.FollowMeeUserName, usr.FollowMeeDeviceID, uss.SyncedThroughDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bcs, err := extractBreadcrumbs(ds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

func getData(followMeeKey string, followMeeUserName string, followMeeDeviceID string, date string) (followMee, error) {
	syncDt, err := time.Parse("20060102", date)
	if err != nil {
		return followMee{}, err
	}

	requestDt := syncDt.Format("2006-01-02")

	url := fmt.Sprintf(
		"http://www.followmee.com/api/tracks.aspx?key=%v&username=%v&output=json&function=daterangefordevice&from=%v&to=%v&deviceid=%v",
		followMeeKey,
		followMeeUserName,
		requestDt,
		requestDt,
		followMeeDeviceID)

	resp, err := http.Get(url)
	if err != nil {
		return followMee{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return followMee{}, fmt.Errorf("error getting data from followMee - %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	data := followMee{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return followMee{}, err
	}

	if data.Error != "" && data.Error != "No data returned for your query" {
		return followMee{}, fmt.Errorf(data.Error)
	}

	return data, nil
}

func extractBreadcrumbs(data followMee) ([]breadcrumb, error) {
	breadcrumbs := []breadcrumb{}

	for _, trackPoint := range data.Data {
		bcTime, err := time.Parse("2006-01-02T15:04:05-07:00", trackPoint.Date)
		if err != nil {
			return []breadcrumb{}, err
		}
		bcTimeString := bcTime.Format("20060102T150405-0700")

		bc := breadcrumb{
			Coordinate: coordinate{
				Lat: trackPoint.Latitude,
				Lon: trackPoint.Longitude,
			},
			Time: bcTimeString,
		}
		breadcrumbs = append(breadcrumbs, bc)
	}

	return breadcrumbs, nil
}
