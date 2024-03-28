package dto

import (
	"net/http"
	"strconv"
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

func TestPostRequest(t *testing.T) {
	var postCases = []struct {
		body     string
		expected bool
		err      string
	}{
		{`{"title": "AD 1", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",}`,
			false, "body should not end with comma"},
		{`{"title": "AD 2", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-24T01:00:00.000Z"}`,
			false, "endAt should be less than startAt"},
		{`{"title": "AD 3", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-32T01:00:00.000Z"}`,
			false, "invalid date"},
		{`{"title": "AD 4", "endAt": "2023-3-31T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z"}`,
			false, "invalid date"},
		{`{"title": "AD 5", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z"}`,
			true, ""},
		{`{"title": "AD 6", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",
			"conditions": {"ageStart": 10, "ageEnd":30}}`,
			true, ""},
		{`{"title": "AD 7", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",
			"conditions": {"gender":["m"]}}`,
			true, ""},
		{`{"title": "AD 8", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",
			"conditions": {"gender":["M"]}}`,
			true, ""},
		{`{"title": "AD 9", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",
			"conditions": {"country":["TW"]}}`,
			true, ""},
		{`{"title": "AD 10", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",
			"conditions": {"country":["tw"]}}`,
			true, ""},
		{`{"title": "AD 11", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",
			"conditions": {"platform":["ios"]}}`,
			true, ""},
		{`{"title": "AD 12", "endAt": "2023-12-23T01:00:00.000Z", "startAt": "2023-12-22T01:00:00.000Z",
			"conditions": {"platform":["IOS"]}}`,
			true, ""},
	}

	t.Run("Test post request", func(t *testing.T) {
		for i, test := range postCases {
			testName := "Test case " + strconv.Itoa(i)
			t.Run(testName, func(t *testing.T) {
				req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/ad", strings.NewReader(test.body))
				_, err := NewRequest(req)

				if err == nil && !test.expected {
					t.Errorf("Test case %d should have error\n %s", i, test.err)
				}

				if err != nil && test.expected {
					t.Errorf("Test case %d should not have error\n %s", i, err)
				}
			})
		}
	})
}

func TestGetRequest(t *testing.T) {
	var getCases = []struct {
		url      string
		expected bool
		err      string
	}{
		{"http://localhost:8080/api/v1/ad?limit=11&age=10&offset=1&offset=10", false, "offset should not be repeated"},
		{"http://localhost:8080/api/v1/ad?limit=101&age=20&offset=10", false, "limit should be less than 100"},
		{"http://localhost:8080/api/v1/ad?limit=-1", false, "limit should be big than 0"},
		{"http://localhost:8080/api/v1/ad?country=tw", false, "country should be iso-3166"},
		{"http://localhost:8080/api/v1/ad?country=usa", false, "country should be iso-3166"},
		{"http://localhost:8080/api/v1/ad?age=1000", false, "age should be less than 100"},
		{"http://localhost:8080/api/v1/ad?limit=10&age=10&offset=1", true, ""},
		{"http://localhost:8080/api/v1/ad?limit=10&age=10&offset=1&country=TW&country=JP&gender=M", true, ""},
		{"http://localhost:8080/api/v1/ad?limit=10&age=10&country=TW&country=JP&gender=M&gender=F", true, ""},
		{"http://localhost:8080/api/v1/ad?limit=10&age=10&offset=1&country=TW&country=JP&gender=M&gender=F&platform=IOS", true, ""},
		{"http://localhost:8080/api/v1/ad?limit=10&age=10&offset=1&country=TW&country=JP&gender=M&gender=F&platform=IOS&platform=Android", true, ""},
		{"http://localhost:8080/api/v1/ad?gender=m&platform=ios", true, ""},
		{"http://localhost:8080/api/v1/ad?gender=M&platform=IOS", true, ""},
		{"http://localhost:8080/api/v1/ad?country=JP,TW", true, ""},
		{"http://localhost:8080/api/v1/ad?gender=M,F", true, ""},
		{"http://localhost:8080/api/v1/ad?platform=IOS,Android", true, ""},
	}

	t.Run("Test get request", func(t *testing.T) {
		for _, test := range getCases {
			testName := test.url
			t.Run(testName, func(t *testing.T) {
				req, _ := http.NewRequest("GET", test.url, nil)
				_, err := NewRequest(req)

				if err == nil && !test.expected {
					t.Errorf("Test case [%s \n] should have error\n %s", testName, test.err)
				}

				if err != nil && test.expected {
					t.Errorf("Test case [%s \n] should not have error\n %s", testName, err)
				}
			})
		}
	})
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
