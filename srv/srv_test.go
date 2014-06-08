package srv

import (
	"testing"

	"github.com/PuerkitoBio/gred/types"
)

func TestSrvGetDB(t *testing.T) {
	s := &server{dbs: make([]DB, maxDBs)}
	// Getting DB 0 should return a non-nil DB
	d0, _ := s.GetDB(0)
	if d0 == nil {
		t.Fatalf("DB at index 0 returned nil")
	}

	// Set a key on d0
	d0.Keys()["a"] = NewKey("a", types.NewString("1"))

	// Get DB 1
	d1, _ := s.GetDB(1)

	// d1 should be != d0
	if d0 == d1 {
		t.Fatalf("DB 0 is the same as DB 1")
	}

	// d1 should not have key "a"
	_, ok := d1.Keys()["a"]
	if ok {
		t.Fatalf("DB 1 has key 'a'")
	}
}
