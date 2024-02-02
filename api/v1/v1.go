package v1

import (
	"fmt"
	"io"
	"net/http"
)

func CreatAD(w http.ResponseWriter, r *http.Request) {
	fmt.Print("v1")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	//read bodt parama
	fmt.Println(string(body))

	//response json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "success"}`))
}

func GetAD(w http.ResponseWriter, r *http.Request) {
	fmt.Print("getAD")

	// Get request body

	//response json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`
	{
		"items": [
		{
		"title": "AD 1",
		"endAt": "2023-12-22T01:00:00.000Z"
		},
		{
		"title": "AD 31",
		"endAt": "2023-12-30T12:00:00.000Z"
		},
		{
		"title": "AD 10",
		"endAt": "2023-12-31T16:00:00.000Z"
		}
	]
	}`))
}
