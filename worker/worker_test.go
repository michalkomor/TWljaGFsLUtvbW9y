package main

import (
	"testing"
)

func TestGetAllURLs(t *testing.T) {
	err := GetAllURLs(&urls)
	if err != nil {
		t.Error(err)
	}
}
