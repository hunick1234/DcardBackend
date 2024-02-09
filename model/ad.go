package model

import (
	"errors"
)

type AD struct {
	Title      string `json:"title"`
	EndAt      string `json:"endAt"`
	StartAt    string `json:"startAt"`
	Timestamp  string `json:"timestamp"`
	Conditions `json:"conditions"`
}

type Conditions struct {
	AgeStart int      `json:"ageStart"`
	AgeEnd   int      `json:"ageEnd"`
	Gender   string   `json:"gender"`
	Country  []string `json:"country"`
	Platform []string `json:"platform"`
}

func (ad *AD) IsValidConditons() error {
	if !ad.Conditions.isValidAge() {
		return errors.New("Invalid age")
	}
	return nil
}

func (ad *Conditions) isValidAge() bool {
	return true
}

func (ad *Conditions) isValidGender() bool {
	return true
}

func (ad *Conditions) isValidCountry() bool {
	return true
}

func (ad *Conditions) isValidPlatform() bool {
	return true
}
