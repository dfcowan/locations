package main

type followMee struct {
	Data []trackPoint `json:"Data"`
}

type trackPoint struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	Date      string  `json:"Date"`
}
