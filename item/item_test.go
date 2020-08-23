package item

import (
	"math/rand"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	collection := New()
	collection.Add(Item{})
	if len(collection.Items) != 1 {
		t.Errorf("Error while adding an url")
	}
}

func TestGetAll(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(10) + 1
	collection := New()
	for i := 0; i < randNum; i++ {
		collection.Add(Item{})
	}
	allItems := collection.GetAll()
	if len(allItems) != randNum {
		t.Errorf("Error while getting the urls")
	}
}

func TestDelete(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(10) + 2
	urls := New()
	for i := 0; i < randNum; i++ {
		urls.Add(Item{})
	}
	urls.Delete(1)
	allURLs := urls.GetAll()
	if len(allURLs) != randNum-1 {
		t.Errorf("Error while deleting an url")
	}
}
