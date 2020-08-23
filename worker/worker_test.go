package main

import (
	"testing"
)

func TestGetAllURLs(t *testing.T) {
	err := GetAllURLs()
	if err != nil {
	 	t.Error(err)
	 }
}
