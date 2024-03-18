package dto

import (
	"net/http"
	"strings"
	"testing"

	"github.com/hunick1234/DcardBackend/model/ad"
)

func TestNewRequest(t *testing.T) {
	// Test case 1: GET request
	getReq, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/ad?limit=11&age=10&offset=1", nil)
	getReqQuery, _ := convertToAdQuery(getReq)
	getReqResult, _ := NewRequest(getReq)
	if adQuery, ok := getReqResult.request.(ad.AdQuery); ok {
		if !compareAdQuery(adQuery, getReqQuery) {
			t.Errorf("SetRequest failed for GET request")
		}
	}

	//test case 2: POST request
	body := `{"title": "AD 1", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",}`
	postReq, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/ad", strings.NewReader(body))
	postReqAd, _ := convertToAd(postReq)
	postReqResult, _ := NewRequest(postReq)
	if ad, ok := postReqResult.request.(ad.AD); ok {
		if !compareAd(ad, postReqAd) {
			t.Errorf("SetRequest failed for POST request")
		}
	}
}

func compareAdQuery(adQuery ad.AdQuery, getReqQuery ad.AdQuery) bool {
	return adQuery.Offset == getReqQuery.Offset &&
		adQuery.Age == getReqQuery.Age &&
		adQuery.Limit == getReqQuery.Limit &&
		compareStringSlice(adQuery.Gender, getReqQuery.Gender) &&
		compareStringSlice(adQuery.Platform, getReqQuery.Platform) &&
		compareStringSlice(adQuery.Country, getReqQuery.Country)
}

func compareAd(ad ad.AD, postReqAd ad.AD) bool {
	return ad.Title == postReqAd.Title &&
		ad.EndAt == postReqAd.EndAt &&
		ad.StartAt == postReqAd.StartAt &&
		ad.Timestamp == postReqAd.Timestamp &&
		ad.StartTimestamp == postReqAd.StartTimestamp &&
		ad.EndTimestamp == postReqAd.EndTimestamp &&
		ad.IsStart == postReqAd.IsStart &&
		ad.IsEnd == postReqAd.IsEnd &&
		compareConditions(ad.Conditions, postReqAd.Conditions)
}

func compareConditions(conditions ad.Conditions, postReqConditions ad.Conditions) bool {
	return conditions.AgeStart == postReqConditions.AgeStart &&
		conditions.AgeEnd == postReqConditions.AgeEnd &&
		compareStringSlice(conditions.Country, postReqConditions.Country) &&
		compareStringSlice(conditions.Gender, postReqConditions.Gender) &&
		compareStringSlice(conditions.Platform, postReqConditions.Platform)
}

func compareStringSlice(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
