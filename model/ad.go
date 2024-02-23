package model

import (
	"fmt"
	"time"

	"github.com/hunick1234/DcardBackend/validation"
	"go.mongodb.org/mongo-driver/bson"
)

type AD struct {
	Title          string     `json:"title" bson:"title"`
	EndAt          string     `json:"endAt" bson:"end_at"`
	StartAt        string     `json:"startAt" bson:"start_at"`
	Timestamp      int64      `json:"timestamp" bson:"timestamp"`
	StartTimestamp int64      `json:"startTimestamp" bson:"start_timestamp"`
	EndTimestamp   int64      `json:"endTimestamp" bson:"end_timestamp"`
	IsStart        bool       `json:"isStart" bson:"is_start"`
	IsEnd          bool       `json:"isEnd" bson:"is_end"`
	Conditions     Conditions `json:"conditions" bson:"conditions"`
}

type ResponseAd struct {
	Title string    `json:"title" bson:"title"`
	EndAt time.Time `json:"endAt" bson:"endAt"`
}

type Conditions struct {
	AgeStart int      `json:"ageStart" bson:"age_start"`
	AgeEnd   int      `json:"ageEnd" bson:"age_end"`
	Gender   []string `json:"gender" bson:"gender"`
	Country  []string `json:"country" bson:"country"`
	Platform []string `json:"platform" bson:"platform"`
}

type AdQuery struct {
	Offset   int      `json:"offset"`
	Limit    int      `json:"limit"`
	Age      int      `json:"age"`
	Gender   []string `json:"gender"`
	Country  []string `json:"country"`
	Platform []string `json:"platform"`
}

var TestValidationAdSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"end_at", "start_at", "timestamp", "conditions", "title"},
		"properties": bson.M{
			"endAt": bson.M{
				"bsonType":    "date",
				"description": "must be a date and is required",
			},
			"conditions": bson.M{
				"bsonType":    "object",
				"description": "must be an object and is required",
				"required":    []string{"age_start", "age_end", "gender", "country", "platform"},
				"properties": bson.M{
					"age_start": bson.M{
						"bsonType":    "int",
						"minimum":     0,
						"maximum":     100,
						"description": "must be an integer and is required",
					},
					"age_end": bson.M{
						"bsonType":    "int",
						"minimum":     0,
						"maximum":     100,
						"description": "must be an integer and is required",
					},
				},
			},
		},
	},
}

func NewAd() AD {
	return AD{
		Title:          "",
		EndAt:          "",
		StartAt:        "",
		Timestamp:      time.Now().Unix(),
		StartTimestamp: 0,
		EndTimestamp:   0,
		IsStart:        false,
		IsEnd:          false,
		Conditions:     NewConditions(),
	}
}

func NewConditions() Conditions {
	return Conditions{
		AgeStart: 0,
		AgeEnd:   100,
		Gender:   []string{},
		Country:  []string{},
		Platform: []string{},
	}
}

func NewAdQuery() AdQuery {
	return AdQuery{
		Offset:   0,
		Limit:    10,
		Age:      0,
		Gender:   []string{},
		Country:  []string{},
		Platform: []string{},
	}
}

func (ad *AD) IsValidAd() bool {

	if !ad.isValidStartAt() || !ad.isValidEndAt() || !ad.isValidAdLiveTime() {
		return false
	}
	return true
}

func (ad *AD) isValidStartAt() bool {
	startAtTime, err := time.Parse(time.RFC3339, ad.StartAt)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	ad.StartTimestamp = startAtTime.Unix()
	return true
}

func (ad *AD) isValidEndAt() bool {

	endAtTime, err := time.Parse(time.RFC3339, ad.EndAt)
	if err != nil {
		return false
	}
	ad.EndTimestamp = endAtTime.Unix()
	return true
}

func (ad *AD) isValidAdLiveTime() bool {
	if ad.StartTimestamp > ad.EndTimestamp {
		return false
	}
	return true
}

func (conditions *Conditions) IsValidConditions() bool {
	if !conditions.isValidAge() || !conditions.isValidGender() || !conditions.isValidCountry() || !conditions.isValidPlatform() {
		return false
	}
	return true
}

func (conditions *Conditions) isValidAge() bool {
	if conditions.AgeStart < 0 || conditions.AgeStart > 100 {
		return false
	}

	if conditions.AgeEnd < 0 || conditions.AgeEnd > 100 {
		return false
	}

	if conditions.AgeStart > conditions.AgeEnd {
		return false
	}
	return true
}

func (conditions *Conditions) isValidGender() bool {
	for _, v := range conditions.Gender {
		if validation.Gender[v] {
			return false
		}
	}
	return true
}

func (conditions *Conditions) isValidCountry() bool {
	for _, v := range conditions.Country {
		if !validation.ISO3166[v] {
			return false
		}
	}
	return true
}

func (conditions *Conditions) isValidPlatform() bool {
	for _, v := range conditions.Platform {
		if validation.Plateform[v] {
			return false
		}
	}
	return true
}

func (query *AdQuery) IsValidAdQuery() bool {
	if !query.isValidOffset() || !query.isValidLimit() || !query.isValidAge() || !query.isValidGender() || !query.isValidCountry() || !query.isValidPlatform() {
		return false
	}
	return true
}

func (query *AdQuery) isValidOffset() bool {
	if query.Offset < 0 {
		return false
	}
	return true
}

func (query *AdQuery) isValidLimit() bool {
	if query.Limit < 0 || query.Limit > 100 {
		return false
	}
	return true
}

func (query *AdQuery) isValidAge() bool {
	if query.Age < 0 || query.Age > 100 {
		return false
	}
	return true
}

func (query *AdQuery) isValidGender() bool {
	for _, v := range query.Gender {
		if validation.Gender[v] {
			return false
		}
	}
	return true
}

func (query *AdQuery) isValidCountry() bool {
	for _, v := range query.Country {
		if !validation.ISO3166[v] {
			return false
		}
	}
	return true
}

func (query *AdQuery) isValidPlatform() bool {
	for _, v := range query.Platform {
		if validation.Plateform[v] {
			return false
		}
	}
	return true
}
