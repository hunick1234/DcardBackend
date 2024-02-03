package model

type AD struct {
	Title     string    `json:"title"`
	EndAt     string    `json:"endAt"`
	StartAt   string    `json:"startAt"`
	Timestamp string    `json:"timestamp"`
	Conditons Conditons `json:"conditons"`
}

type Conditons struct {
	AgeStart int      `json:"ageStart"`
	AgeEnd   int      `json:"ageEnd"`
	Country  []string `json:"country"`
	Platform []string `json:"platform"`
}
