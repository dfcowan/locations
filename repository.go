package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var (
	connStr string
)

func migrate() error {
	fmt.Println("Migrating database")

	connStr = os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://locations:locations@localhost/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query("SELECT u_followMee_key FROM users WHERE user_id = 0")
	if err != nil {
		if err.Error() == "pq: relation \"users\" does not exist" {
			fmt.Println("Creating users table")
			_, err := db.Exec(
				`create table if not exists users (
					user_id bigint primary key,
					u_followMee_key char(64) not null default '',
					u_followMee_username char(64) not null default '',
					u_followMee_deviceid char(64) not null default ''
				)`)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Created users table")
		} else if err.Error() == "pq: column \"u_followmee_key\" does not exist" {
			fmt.Println("Altering users table")
			_, err := db.Exec(
				`alter table users 
					add column u_followMee_key char(64) not null default '',
					add column u_followMee_username char(64) not null default '',
					add column u_followMee_deviceid char(64) not null default ''
				`)
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

func createUser(userID string, followMeeKey string, followMeeUserName string, followMeeDeviceID string, firstDate string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO users(user_id, u_followMee_key, u_followMee_username, u_followMee_deviceid) 
		values($1, $2, $3)`,
		userID,
		followMeeKey,
		followMeeUserName,
		followMeeDeviceID)
	if err != nil {
		return err
	}
	fmt.Println("inserted to users")

	_, err = db.Exec(
		`INSERT INTO users_sync_status(user_id, uss_start_date, uss_synced_through_date) 
		values($1, $2, $3)`,
		userID,
		firstDate,
		firstDate)
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

	fmt.Println(fmt.Sprintf("updated synced through date %v", syncedThroughDate))

	return nil
}

func loadUser(userID int) (user, userSyncStatus, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return user{}, userSyncStatus{}, err
	}
	defer db.Close()

	row := db.QueryRow(
		`select user_id, u_followMee_key, u_followMee_username, u_followMee_deviceid
		from users
		where user_id = $1`,
		userID)

	u := user{}

	err = row.Scan(&u.UserID, &u.FollowMeeKey, &u.FollowMeeUserName, &u.FollowMeeDeviceID)
	if err == sql.ErrNoRows {
		return user{}, userSyncStatus{}, errors.New("user not found")
	} else if err != nil {
		return user{}, userSyncStatus{}, err
	}
	u.FollowMeeKey = strings.Trim(u.FollowMeeKey, " ")
	u.FollowMeeUserName = strings.Trim(u.FollowMeeUserName, " ")
	u.FollowMeeDeviceID = strings.Trim(u.FollowMeeDeviceID, " ")

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

	if uss.SyncedThroughDate < uss.StartDate {
		uss.SyncedThroughDate = uss.StartDate
	}

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

	fmt.Println(fmt.Sprintf("saved breadcrumbs %v", breadcrumbs))

	return nil
}

func deleteBreadCrumbs(userID int) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	sdp := "2019-05-07T000000-0000"
	edp := "2019-05-08T235959-0000"

	_, err = db.Exec(
		`delete
		from breadcrumbs
		where user_id = $1
		  and bc_time between $2 and $3`,
		userID,
		sdp,
		edp)
	if err != nil {
		return err
	}

	return nil
}

func loadCounts(userID int, startDate string, endDate string) ([]coordinateCount, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return []coordinateCount{}, err
	}
	defer db.Close()

	var sdp, edp string
	if startDate == "" {
		sdp = "00000000-T000000-0000"
	} else {
		sdp = startDate + "T000000-0000"
	}
	if endDate == "" {
		edp = "99991231-T000000-0000"
	} else {
		edp = endDate + "T235959-0000"
	}

	rows, err := db.Query(
		`select bc_lat, bc_lon, count(*)
		from breadcrumbs
		where user_id = $1
		  and bc_time between $2 and $3
		group by bc_lat, bc_lon`,
		userID,
		sdp,
		edp)
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
