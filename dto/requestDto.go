package dto

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hunick1234/DcardBackend/model/ad"
)

type Request struct {
	request any
}

func NewRequest(r *http.Request) (Request, error) {
	switch r.Method {
	case "GET":
		req, err := convertToAdQuery(r)
		if err != nil {
			return Request{}, err
		}
		return Request{request: req}, nil

	case "POST":
		req, err := convertToAd(r)
		if err != nil {
			return Request{}, err
		}
		return Request{request: req}, nil
	default:
		return Request{}, errors.New("invalid request")
	}
}

func (r *Request) SetRequest(request any) {
	r.request = request
}

func (r *Request) GetRequest() any {
	return r.request
}

func (r *Request) GetRequestAd() ad.AD {
	if ad, ok := r.request.(ad.AD); ok {
		return ad
	}
	return ad.AD{}
}

func (r *Request) GetRequestAdQuery() ad.AdQuery {
	if adQuery, ok := r.request.(ad.AdQuery); ok {
		return adQuery
	} else {
		return ad.AdQuery{}
	}
}

func convertToAdQuery(r *http.Request) (ad.AdQuery, error) {
	query := ad.DefaultAdQuery()

	offset := r.URL.Query()["offset"]
	if result, ok := query2Int(offset); ok {
		query.Offset = result
	} else {
		return ad.AdQuery{}, errors.New("wrong query")
	}

	limit := r.URL.Query()["limit"]
	if result, ok := query2Int(limit); ok {
		query.Limit = result
	} else {
		return ad.AdQuery{}, errors.New("wrong query")
	}

	age := r.URL.Query()["age"]
	if result, ok := query2Int(age); ok {
		query.Age = result
	} else {
		return ad.AdQuery{}, errors.New("wrong query")
	}

	gender := r.URL.Query()["gender"]
	query.Gender = query2Arr(gender)
	country := r.URL.Query()["country"]
	query.Country = query2Arr(country)
	platform := r.URL.Query()["platform"]
	query.Platform = query2Arr(platform)

	if !query.IsValidAdQuery() {
		return ad.AdQuery{}, errors.New("invalid query vaule")
	}
	return query, nil
}

func convertToAd(r *http.Request) (ad.AD, error) {
	var ad ad.AD
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ad, err
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &ad)
	if err != nil {
		return ad, err
	}

	if !ad.IsValidAd() {
		return ad, errors.New("invalid ad")
	}
	if !ad.Conditions.IsValidConditions() {
		return ad, errors.New("invalid conditions")
	}
	return ad, nil
}

func query2Int(q []string) (int, bool) {
	if len(q) == 0 {
		return 0, true
	}
	if len(q) > 1 {
		return 0, false
	}
	if result, err := strconv.Atoi(q[0]); err == nil {
		return result, true
	}
	return 0, false
}

// 暫不支援混用.  ?country=JP,TW&country=US....
func query2Arr(s []string) []string {
	if len(s) == 0 {
		return []string{}
	}
	if len(s) == 1 {
		if strings.Contains(s[0], ",") {
			return splitByComma(s[0])
		}
	}
	return s
}

func splitByComma(s string) []string {
	return strings.Split(s, ",")
}
