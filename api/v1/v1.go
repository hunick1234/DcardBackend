package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/hunick1234/DcardBackend/model/ad"
)

func CreatAD(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	//read body parama
	fmt.Println(string(body))
	storeItem := ad.NewAd()
	err = json.Unmarshal(body, &storeItem)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	//vaildate ad
	if !storeItem.IsValidAd() {
		http.Error(w, "Invalid ad", http.StatusBadRequest)
		return
	}
	//store ad
	err = ad.DeafultAdService().Store(&storeItem)

	if err != nil {
		http.Error(w, "Failed to store ad", http.StatusBadRequest)
		return
	}

	//response json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "success"}`))
}

func GetAD(w http.ResponseWriter, r *http.Request) {
	query := ad.NewAdQuery()

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
		http.Error(w, "Invalid query", http.StatusBadRequest)
		return
	}
	//get ad
	rep, err := ad.DeafultAdService().FindByFilter(context.TODO(), query)
	if err != nil {
		http.Error(w, "Failed to get ad", http.StatusBadRequest)
		return
	}

	//response json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rep)

}
