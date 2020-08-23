package main

import (
	"GWP/item"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

//type TestServer struct{}

func TestInitialState(t *testing.T) {

	// //router := chi.NewRouter()
	// //req := httptest.NewRequest(http.MethodGet, "/api/fetcher", nil)
	// //recorder := httptest.NewRecorder()
	// //srv := httptest.NewServer(router)
	// resp, err := http.Get("http://localhost:8080/api/fetcher")
	// if err != nil {
	// 	t.Error(err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Error(err)
	// }
	//
	// //GetAllURLs(w, r)
	// expectedResponse := "[]"
	// //t.Errorf("Expects %s but received %s", expectedResponse, body)
	// if string(body) != expectedResponse {
	// 	t.Errorf("Expects %s but received %s ", expectedResponse, string(body))
	// }
}

func TestGetWrongID(t *testing.T) {

	router := chi.NewRouter()

	r := httptest.NewRequest(http.MethodGet, "/api/fetcher/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	expectedHeader := http.StatusNotFound
	if w.Code != expectedHeader {
		t.Errorf("Expects %d but received %d", expectedHeader, w.Code)
	}
}

func TestPostGetAndDeleteURL(t *testing.T) {

	//ts := httptest.NewServer(r)
	httptest.NewServer(r)
	//defer ts.Close()
	//create new json with url to post
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
	//router.ServeHTTP(w, r)
	AddURL(w1, r1)
	expectedResponse := "{\"id\": 1}"
	if w1.Body.String() != expectedResponse {
		t.Errorf("Expects %s but received %s", expectedResponse, w1.Body.String())
	}
	//get URL
	r2 := httptest.NewRequest(http.MethodGet, "/api/fetcher", bytes.NewReader(nil))
	w2 := httptest.NewRecorder()
	GetAllURLs(w2, r2)
	jsonbytes, _ := json.Marshal(collection.Items)
	expectedResponse = string(jsonbytes)
	if w2.Body.String() != expectedResponse {
		t.Errorf("Expects %s but received %s", expectedResponse, w2.Body.String())
	}
	//delete request
	// r3 := httptest.NewRequest(http.MethodDelete, "/api/fetcher/1", nil)
	// w3 := httptest.NewRecorder()
	// //
	// DeleteURL(w3, r3)
	// //
	// expectedResponse = "item 1 deleted"
	// // if w3.Body.String() != expectedResponse {
	// t.Errorf("Expects %s but received %s", expectedResponse, w3.Body.String())
	// }
}
