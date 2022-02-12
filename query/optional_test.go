package query_test

import (
	"errors"
	"testing"

	"github.com/gopherd/doge/query"
)

func TestOr(t *testing.T) {
	if v := query.Or("x", "y"); v != "x" {
		t.Fatalf("want %q, but got %q", "x", v)
	}
	if v := query.Or("", "y"); v != "y" {
		t.Fatalf("want %q, but got %q", "y", v)
	}
	if v := query.Or(1, 2); v != 1 {
		t.Fatalf("want 1, but got %d", v)
	}
	if v := query.Or(0, 2); v != 2 {
		t.Fatalf("want 1, but got %d", v)
	}
	if v := query.Or(true, true); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := query.Or(true, false); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := query.Or(false, true); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := query.Or(false, false); v != false {
		t.Fatalf("want false, but got %v", v)
	}

	var errTest1 = errors.New("test1")
	var errTest2 = errors.New("test2")
	if v := query.Or(errTest1, nil); v != errTest1 {
		t.Fatalf("want %v, but got %v", errTest1, v)
	}
	if v := query.Or(nil, errTest1); v != errTest1 {
		t.Fatalf("want %v, but got %v", errTest1, v)
	}
	if v := query.Or(errTest1, errTest2); v != errTest1 {
		t.Fatalf("want %v, but got %v", errTest1, v)
	}
}
