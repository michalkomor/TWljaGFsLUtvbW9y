package main

import (
	"GWP/item"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	//"github.com/go-chi/chi"
)

func TestInitialState(t *testing.T) {

	r1 := httptest.NewRequest(http.MethodGet, "/api/fetcher", nil)
	w1 := httptest.NewRecorder()
	GetAllURLs(w1, r1)

	bytes, _ := json.Marshal(collection.Items)
	expectedResponse := string(bytes)
	if w1.Body.String() != expectedResponse {
		t.Errorf("Expects %s but received %s ", expectedResponse, w1.Body.String())
	}
}


func TestPostAndGetURL(t *testing.T) {

	httptest.NewServer(r)


	testURL := item.Item{
		ID:       1,
		URL:      "https://httpbin.org/range/1",
		Interval: 5,
	}
	var collection = item.New()
	collection.Add(testURL)
	bts, err := json.Marshal(testURL)
	if err != nil {
		t.Error(err)
	}
	//post URL
	r1 := httptest.NewRequest(http.MethodPost, "/api/fetcher", bytes.NewReader(bts))
	w1 := httptest.NewRecorder()

	AddURL(w1, r1)
	expectedResponse := "{\"id\": 1}"
	if w1.Body.String() != expectedResponse {
		t.Errorf("Expects %s but received %s", expectedResponse, w1.Body.String())
	}
	//get URL
	r2 := httptest.NewRequest(http.MethodGet, "/api/fetcher", bytes.NewReader(nil))
	w2 := httptest.NewRecorder()
	GetAllURLs(w2, r2)
	jsonbytes, err := json.Marshal(collection.Items)
	if err != nil {
		t.Error(err)
	}
	expectedResponse = string(jsonbytes)
	if w2.Body.String() != expectedResponse {
		t.Errorf("Expects %s but received %s", expectedResponse, w2.Body.String())
	}
	//delete request
	r3 := httptest.NewRequest(http.MethodDelete, "/api/fetcher/1", nil)
	w3 := httptest.NewRecorder()
	// //
	DeleteURL(w3, r3)
	// //
	expectedResponse = "item 1 deleted"
	// // if w3.Body.String() != expectedResponse {
	t.Errorf("Expects %s but received %s", expectedResponse, w3.Body.String())
	// }
}


func TestGetWrongID(t *testing.T) {
	httptest.NewServer(r)
	r1 := httptest.NewRequest(http.MethodGet, "/api/fetcher/1", nil)
	w1 := httptest.NewRecorder()
	GetURL(w1, r1)
	expectedHeader := http.StatusNotFound
	if w1.Code != expectedHeader {
		t.Errorf("Expects %d but received %d", expectedHeader, w1.Code)
	}
}
