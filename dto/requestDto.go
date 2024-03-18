package dto

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

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

	offsetStr := r.URL.Query().Get("offset")
	if offset, err := strconv.Atoi(offsetStr); err == nil {
		query.Offset = offset
	}
	limitStr := r.URL.Query().Get("limit")
	if limit, err := strconv.Atoi(limitStr); err == nil {
		query.Limit = limit
	}
	ageStr := r.URL.Query().Get("age")
	if age, err := strconv.Atoi(ageStr); err == nil {
		query.Age = age
	}
	gender := r.URL.Query()["gender"]
	query.Gender = gender
	country := r.URL.Query()["country"]
	query.Country = country
	platform := r.URL.Query()["platform"]
	query.Platform = platform

	if !query.IsValidAdQuery() {
		return ad.AdQuery{}, errors.New("invalid query")
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
