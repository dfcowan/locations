package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	connStr string
)

func migrate() error {
	fmt.Println("Migrating database")

	connStr = os.Getenv("CONNSTRING")
	if connStr == "" {
		connStr = "postgres://locations:locations@localhost/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query("SELECT user_id FROM users WHERE user_id = 0")
	if err != nil {
		if err.Error() == "pq: relation \"users\" does not exist" {
			fmt.Println("Creating users table")
			_, err := db.Exec(
				`create table if not exists users (
					user_id bigint primary key,
					u_access_token char(64) not null,
					u_refresh_token char(64) not null
				)`)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Created users table")
		} else {
			log.Fatal(err)
		}
	}

	_, err = db.Query("SELECT user_id FROM users_sync_status WHERE user_id = 0")
	if err != nil {
		if err.Error() == "pq: relation \"users_sync_status\" does not exist" {
			fmt.Println("Creating users_sync_status table")
			_, err := db.Exec(
				`create table if not exists users_sync_status (
					user_id bigint primary key,
					uss_start_date char(64) not null,
					uss_synced_through_date char(64) not null
				)`)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Created users_sync_status table")
		} else {
			log.Fatal(err)
		}
	}

	_, err = db.Query("SELECT user_id FROM breadcrumbs WHERE user_id = 0")
	if err != nil {
		if err.Error() == "pq: relation \"breadcrumbs\" does not exist" {
			fmt.Println("Creating breadcrumbs table")
			_, err := db.Exec(
				`create table if not exists breadcrumbs (
					user_id bigint not null,
					bc_lat numeric(9, 6) not null,
					bc_lon numeric(9, 6) not null,
					bc_time char(32) not null
				)`)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Created breadcrumbs table")
		} else {
			log.Fatal(err)
		}
	}

	_, err = db.Query("SELECT user_id FROM coordinate_counts WHERE user_id = 0")
	if err != nil {
		if err.Error() == "pq: relation \"coordinate_counts\" does not exist" {
			fmt.Println("Creating coordinate_counts table")
			_, err := db.Exec(
				`create table if not exists coordinate_counts (
					user_id bigint not null,
					cc_lat numeric(9, 6) not null,
					cc_lon numeric(9, 6) not null,
					cc_count bigint not null
				)`)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Created coordinate_counts table")
		} else {
			log.Fatal(err)
		}
	}

	return nil
}

func createUser(token accessToken, firstDate string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO users(user_id, u_access_token, u_refresh_token) 
		values($1, $2, $3)`,
		token.UserID,
		token.AccessToken,
		token.RefreshToken)
	if err != nil {
		return err
	}
	fmt.Println("inserted to users")

	dt, err := time.Parse("20060102", firstDate)
	if err != nil {
		return err
	}
	dt = dt.AddDate(0, 0, -1)
	syncedDate := dt.Format("20060102")

	_, err = db.Exec(
		`INSERT INTO users_sync_status(user_id, uss_start_date, uss_synced_through_date) 
		values($1, $2, $3)`,
		token.UserID,
		firstDate,
		syncedDate)
	if err != nil {
		return err
	}
	fmt.Println("inserted to users_sync_status")

	return nil
}

func saveUserSyncedThroughDate(userID int, syncedThroughDate string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`update users_sync_status 
		set uss_synced_through_date = $1
		where user_id = $2`,
		syncedThroughDate,
		userID)
	if err != nil {
		return err
	}
	fmt.Println("updated synced through date")

	return nil
}

func loadUser(userID int) (user, userSyncStatus, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return user{}, userSyncStatus{}, err
	}
	defer db.Close()

	row := db.QueryRow(
		`select user_id, u_access_token, u_refresh_token
		from users
		where user_id = $1`,
		userID)

	u := user{}

	err = row.Scan(&u.UserID, &u.AccessToken, &u.RefreshToken)
	if err == sql.ErrNoRows {
		return user{}, userSyncStatus{}, errors.New("user not found")
	} else if err != nil {
		return user{}, userSyncStatus{}, err
	}

	row = db.QueryRow(
		`select user_id, uss_start_date, uss_synced_through_date
		from users_sync_status
		where user_id = $1`,
		userID)

	uss := userSyncStatus{}

	err = row.Scan(&uss.UserID, &uss.StartDate, &uss.SyncedThroughDate)
	if err == sql.ErrNoRows {
		return user{}, userSyncStatus{}, errors.New("user sync status not found")
	} else if err != nil {
		return user{}, userSyncStatus{}, err
	}

	uss.StartDate = strings.Trim(uss.StartDate, " ")
	uss.SyncedThroughDate = strings.Trim(uss.SyncedThroughDate, " ")

	return u, uss, nil
}

func saveBreadcrumbs(userID int, breadcrumbs []breadcrumb) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, bc := range breadcrumbs {

		_, err = tx.Exec(
			`INSERT INTO breadcrumbs(user_id, bc_lat, bc_lon, bc_time) 
			values($1, $2, $3, $4)`,
			userID,
			bc.Coordinate.Lat,
			bc.Coordinate.Lon,
			bc.Time)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Println("saved breadcrumbs")

	return nil
}

func loadCounts(userID int) ([]coordinateCount, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return []coordinateCount{}, err
	}
	defer db.Close()

	rows, err := db.Query(
		`select bc_lat, bc_lon, count(*)
		from breadcrumbs
		where user_id = $1
		group by bc_lat, bc_lon`,
		userID)
	if err != nil {
		return []coordinateCount{}, err
	}

	defer rows.Close()

	counts := []coordinateCount{}
	for rows.Next() {
		cc := coordinateCount{}

		err = rows.Scan(&cc.Coordinate.Lat, &cc.Coordinate.Lon, &cc.Count)
		if err != nil {
			return []coordinateCount{}, err
		}

		counts = append(counts, cc)
	}

	return counts, nil
}