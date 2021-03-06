package main

type user struct {
	UserID            int
	FollowMeeKey      string
	FollowMeeUserName string
	FollowMeeDeviceID string
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
	P coordinate
	C int
}

type status struct {
	UserCount       int
	BreadcrumbCount int
}
