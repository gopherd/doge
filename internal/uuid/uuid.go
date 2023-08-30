// package uuid implements an UUID allocator, it's used to internal only.
//
// NOTE: This package was copied from @https://github.com/google/uuid without unused parts.
package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC 4122.
type UUID [16]byte

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (uuid UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// New returns a Random (Version 4) UUID.
func New() (UUID, error) {
	return NewFromReader(rand.Reader)
}

// NewString returns a Random (Version 4) UUID string
func NewString() string {
	uuid, err := New()
	if err != nil {
		panic(err)
	}
	return uuid.String()
}

// NewFromReader returns a UUID based on bytes read from a given io.Reader.
func NewFromReader(r io.Reader) (UUID, error) {
	var uuid UUID
	_, err := io.ReadFull(r, uuid[:])
	if err != nil {
		return UUID{}, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid, nil
}
