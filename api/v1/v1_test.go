package v1_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	v1 "github.com/hunick1234/DcardBackend/api/v1"
)

func TestCreateAD(t *testing.T) {
	// Create a new HTTP request
	body := strings.NewReader(`{"title": "AD 1", "endAt": "2023-12-22T01:00:00.000Z"}`)
	req, err := http.NewRequest("POST", "post/v1/ads", body)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(v1.CreatAD)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	expected := `{"message": "success"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetAD(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "get/v1/ads", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(v1.GetAD)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	expected := `
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
	}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
