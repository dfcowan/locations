package main

type user struct {
	UserID       int
	AccessToken  string
	RefreshToken string
}

type userSyncStatus struct {
	UserID            int
	StartDate         string
	SyncedThroughDate string
}

type breadcrumb struct {
	Coordinate coordinate
	Time       string
}

type coordinate struct {
	Lat float64
	Lon float64
}

type coordinateCount struct {
	Coordinate coordinate
	Count      int
}
