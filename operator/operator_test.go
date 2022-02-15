package operator_test

import (
	"errors"
	"testing"

	"github.com/gopherd/doge/operator"
)

func TestOr(t *testing.T) {
	if v := operator.Or("x", "y"); v != "x" {
		t.Fatalf("want %q, but got %q", "x", v)
	}
	if v := operator.Or("", "y"); v != "y" {
		t.Fatalf("want %q, but got %q", "y", v)
	}
	if v := operator.Or(1, 2); v != 1 {
		t.Fatalf("want 1, but got %d", v)
	}
	if v := operator.Or(0, 2); v != 2 {
		t.Fatalf("want 1, but got %d", v)
	}
	if v := operator.Or(true, true); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := operator.Or(true, false); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := operator.Or(false, true); v != true {
		t.Fatalf("want true, but got %v", v)
	}
	if v := operator.Or(false, false); v != false {
		t.Fatalf("want false, but got %v", v)
	}

	var errTest1 = errors.New("test1")
	var errTest2 = errors.New("test2")
	if v := operator.Or(errTest1, nil); v != errTest1 {
		t.Fatalf("want %v, but got %v", errTest1, v)
	}
	if v := operator.Or(nil, errTest1); v != errTest1 {
		t.Fatalf("want %v, but got %v", errTest1, v)
	}
	if v := operator.Or(errTest1, errTest2); v != errTest1 {
		t.Fatalf("want %v, but got %v", errTest1, v)
	}
}
