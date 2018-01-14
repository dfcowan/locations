package main

type accessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	UserID       int    `json:"user_id"`
	Error        string `json:"error"`
}

type dailyStoryline struct {
	Date     string             `json:"date"`
	Segments []storylineSegment `json:"segments"`
}

type storylineSegment struct {
	Type       string     `json:"type"`
	Activities []activity `json:"activities"`
	LastUpdate string     `json:"lastUpdate"`
}

type activity struct {
	Activity    string       `json:"activity"`
	TrackPoints []trackPoint `json:"trackPoints"`
}

type trackPoint struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Time string  `json:"time"`
}