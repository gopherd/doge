package proto_test

import (
	"testing"

	. "github.com/gopherd/doge/proto"
	"github.com/gopherd/doge/proto/testdata"
)

func TestEncodeDecode(t *testing.T) {
	m := &testdata.Testdata{
		Version: 14,
		Id:      23,
		Name:    "hello",
	}
	buf, err := Encode(m, 0)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	t.Logf("encoded: %v", buf)

	n, m0, err := Decode(buf)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("remain %d bytes unread", len(buf)-n)
	}
	m1, ok := m0.(*testdata.Testdata)
	if !ok {
		t.Fatalf("unexpected message type: %T", m0)
	}
	if m.Version != m1.Version || m.Id != m1.Id || m.Name != m1.Name {
		t.Fatalf("decoded message mismatched: %v vs %v", m, m1)
	}
}
