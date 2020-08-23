package main

import (
	"GWP/item"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var maxSize int64 = (1 << 20) //1MB
var index int = 0             //initial index value

var collection = item.New() //slice to store URLs

var r = chi.NewRouter()

//AddURL - Add URL to slice of URLs
func AddURL(w http.ResponseWriter, r *http.Request) {
	//checks if request size does not exceed maxSize (1MB)
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	//new decoder that reads from r
	dec := json.NewDecoder(r.Body)
	//checks if request body contains only fields that are in the structure
	dec.DisallowUnknownFields()
	//decode request in json format and store it in the url struct
	var newItem item.Item
	err := dec.Decode(&newItem)
	//check various errors
	if err != nil {
		handleError(err, w)
		return
	}
	//there was no error - add URL to slice of URLs
	index++
	newItem.ID = index
	collection.Add(newItem)
	w.WriteHeader(200)
	msg := fmt.Sprintf("{\"id\": %d}", newItem.ID)
	w.Write([]byte(msg))
	return
}

//GetAllURLs - list all URLs in the slice
func GetAllURLs(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal(collection.Items)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error occured while marshaling to JSON"))
		return
	}
	w.Write(jsonBytes)
	return
}

//GetURL - sends json with one URL
func GetURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(urlID)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	for _, currURL := range collection.GetAll() {
		if currURL.ID == ID {
			//url = currURL
			jsonBytes, err := json.Marshal(currURL)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Error occured while marshaling to JSON"))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(jsonBytes))
			return
		}
	}
	w.WriteHeader(404)
	return
}

//DeleteURL - delete the URL
func DeleteURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(urlID)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	for i, currURL := range collection.GetAll() {
		if currURL.ID == ID {
			collection.Delete(i)
			msg := fmt.Sprintf("item %d deleted", ID)
			w.WriteHeader(200)
			w.Write([]byte(msg))
			return
		}
	}
	//id has not been found - response status 404
	w.WriteHeader(404)
	return
}

//UpdateURL - update the URL;
func UpdateURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(urlID)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	for i, currURL := range collection.GetAll() {
		if currURL.ID == ID {
			//checks if request size does not exceed maxSize (1MB)
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			//new decoder that reads from r
			dec := json.NewDecoder(r.Body)
			//checks if request body contains only fields that are in the structure
			dec.DisallowUnknownFields()
			//decodes request in json format and store it in the url struct
			var newItem item.Item
			err := dec.Decode(&newItem)
			//checks various errors
			if err != nil {
				handleError(err, w)
				return
			}
			//there was no error - update URL
			collection.Items[i].URL = newItem.URL
			collection.Items[i].Interval = newItem.Interval
			//delete the history of old url
			collection.Items[i].History = make([]item.History, 0)
			w.WriteHeader(200)
			msg := fmt.Sprintf("{\"id\": %d} updated", collection.Items[i].ID)
			w.Write([]byte(msg))
			w.WriteHeader(200)
			return
		}
	}
	//id has not been found - response status 404
	w.WriteHeader(404)
	return
}

//GetURLHistory - sends json with url request history
func GetURLHistory(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(urlID)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	for _, currURL := range collection.GetAll() {
		if currURL.ID == ID {
			jsonBytes, err := json.Marshal(currURL.History)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Error occured while marshaling to JSON"))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(jsonBytes))
			return
		}
	}
	//id has not been found - response status 404
	w.WriteHeader(404)
	return
}

func AddURLHistory(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(urlID)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	for i, currURL := range collection.GetAll() {
		if currURL.ID == ID {
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			//new decoder that reads from r
			dec := json.NewDecoder(r.Body)
			//checks if request body contains only fields that are in the structure
			dec.DisallowUnknownFields()
			//decodes request in json format and store it in the url struct
			var newHistory item.History
			err := dec.Decode(&newHistory)
			if err != nil {
				handleError(err, w)
				return
			}
			collection.Items[i].History = append(collection.Items[i].History, newHistory)
			w.WriteHeader(200)
			return
		}
	}
}

func handleError(err error, w http.ResponseWriter) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	switch {
	// Catch the error caused by the request body being too large.
	case err.Error() == "http: request body too large":
		msg := "Request body must not be larger than 1MB"
		http.Error(w, msg, http.StatusRequestEntityTooLarge)
	//SyntaxError
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON")
		http.Error(w, msg, http.StatusBadRequest)

	//Unexpected end of file
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := fmt.Sprintf("Request body contains badly-formed JSON")
		http.Error(w, msg, http.StatusBadRequest)

	// Catch any type errors, like trying to assign a string in the
	// JSON request body to a int field in the struct
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		http.Error(w, msg, http.StatusBadRequest)

	// Catch the error caused by extra unexpected fields in the request body.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		http.Error(w, msg, http.StatusBadRequest)

	// An io.EOF error is returned by Decode() if the request body is empty.
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		http.Error(w, msg, http.StatusBadRequest)

	// Server Error response.
	default:
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {

	// Logs the start and end of each request with the elapsed processing time
	r.Use(middleware.Logger)
	// Injects a request ID into the context of each request
	r.Use(middleware.RequestID)
	//Parse extension from url and put it on request context
	r.Use(middleware.URLFormat)
	//Gracefully absorb panics and prints the stack trace
	r.Use(middleware.Recoverer)

	r.Get("/api/fetcher", GetAllURLs) //Read
	r.Post("/api/fetcher", AddURL)    //Create

	r.Route("/api/fetcher/{id}", func(r chi.Router) {
		r.Get("/", GetURL)                //Read
		r.Delete("/", DeleteURL)          //Delete
		r.Get("/history", GetURLHistory)  //Read
		r.Put("/", UpdateURL)             //Update
		r.Post("/history", AddURLHistory) //Create
	})
	http.ListenAndServe(":8080", r)
}
