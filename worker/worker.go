package main

import (
	"GWP/item"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var requestTimeout = 5 * time.Second //5 seconds
var collection = item.New()          //slice to store URLs

//GetAllURLs - sends json file with all exisitng url
func GetAllURLs() error {
	r, err := http.Get("http://localhost:8080/api/fetcher")
	if err != nil {
		log.Println(err)
		return err
	}
	defer r.Body.Close()
	//new decoder that reads from r
	dec := json.NewDecoder(r.Body)
	//checks if request body contains only fields that are in the structure
	dec.DisallowUnknownFields()
	//decode request in json format and store it in the url struct
	err = dec.Decode(&collection.Items)
	//check various errors
	if err != nil {
		handleError(err)
		return err
	}
	return nil
}

func handleError(err error) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	switch {
	// Catch the error caused by the request body being too large.
	case err.Error() == "http: request body too large":
		msg := "Request body must not be larger than 1MB"
		log.Println(msg)
	//SyntaxError
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON")
		log.Println(msg)

	//Unexpected end of file
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := fmt.Sprintf("Request body contains badly-formed JSON")
		log.Println(msg)

	// Catch any type errors, like trying to assign a string in the
	// JSON request body to a int field in the struct
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		log.Println(msg)

	// Catch the error caused by extra unexpected fields in the request body.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		log.Println(msg)

	// An io.EOF error is returned by Decode() if the request body is empty.
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		log.Println(msg)

	// Server Error response.
	default:
		log.Println(err.Error())
	}
}

func main() {
	var wg sync.WaitGroup
	GetAllURLs()
	urlsLength := len(collection.GetAll())
	wg.Add(urlsLength)
	for i, _ := range collection.GetAll() {
		go func(i int) {
			defer wg.Done()
			var msg string
			var hist item.History
			client := http.Client{
				Timeout: 5 * time.Second,
			}
			ticker := time.NewTicker(time.Duration((collection.Items[i].Interval)) * time.Second)
			for {
				select {
				case <-ticker.C:
					{
						urlItem := collection.Items[i]
						//hist = make([]url.History, 0)
						start := time.Now()
						r, err := client.Get(urlItem.URL)
						secs := time.Since(start).Seconds()
						if err != nil {
							msg = err.Error()
							hist = item.History{Response: "", Duration: secs, CreatedAt: time.Now().Unix()}
							collection.Items[i].History = append(collection.Items[i].History, hist)
						} else {
							defer r.Body.Close()
							data, err := ioutil.ReadAll(r.Body)
							if err != nil {
								msg = err.Error()
								hist = item.History{Response: "", Duration: secs, CreatedAt: time.Now().Unix()}
								collection.Items[i].History = append(collection.Items[i].History, hist)
							} else {
								hist = item.History{Response: string(data), Duration: secs, CreatedAt: time.Now().Unix()}
								collection.Items[i].History = append(collection.Items[i].History, hist)
								msg = string(data)
							}
						}
						jsonBytes, err := json.Marshal(hist)
						if err != nil {
							log.Println(err.Error())
						}
						temp := fmt.Sprintf("%d", collection.Items[i].ID)
						postURL := fmt.Sprintf("http://localhost:8080/api/fetcher/" + temp + "/history")
						_, err = http.Post(postURL, "application/json", bytes.NewBuffer(jsonBytes))
						if err != nil {
							log.Println(err.Error())
						}
						fmt.Println(msg)
					}
				}
			}
		}(i)
	}
	wg.Wait()
}
