package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

	r2 := httptest.NewRequest(http.MethodGet, "/api/fetcher/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r2)

	//GetURL(w, r2)
	expectedHeader := http.StatusNotFound
	if w.Code != expectedHeader {
		t.Errorf("Expects %d but received %d", expectedHeader, w.Code)
	}
}

func TestPostGetAndDeleteURL(t *testing.T) {
	//create new json with url to post
	testURL := url{
		ID:       1,
		URL:      "https://httpbin.org/range/1",
		Interval: 5,
	}
	bts, err := json.Marshal(testURL)
	if err != nil {
		t.Error(err)
	}
	//post URL
	r := httptest.NewRequest(http.MethodPost, "/api/fetcher", bytes.NewReader(bts))
	w := httptest.NewRecorder()
	AddURL(w, r)
	expectedResponse := "{\"id\": 1}"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expects %s but received %s", expectedResponse, w.Body.String())
	}
	//get URL

	//delete request
	// r = httptest.NewRequest(http.MethodDelete, "/api/fetcher/1", nil)
	// w = httptest.NewRecorder()
	//
	// DeleteURL(w, r)
	//
	// expectedResponse = "item 1 deleted"
	// if w.Body.String() != expectedResponse {
	// 	t.Errorf("Expects %s but received %s", expectedResponse, w.Body.String())
	// }
}
