package routercache

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"time"
)

// Response is the cached response data structure.
type Response struct {
	// Value is the cached response value.
	Value []byte

	// Header is the cached response header.
	Header http.Header

	// StatusCode is the cached response status code.
	StatusCode int

	// Expiration is the cached response expiration date.
	Expiration time.Time

	// LastAccess is the last date a cached response was accessed.
	// Used by LRU and MRU algorithms.
	LastAccess time.Time

	// Frequency is the count of times a cached response is accessed.
	// Used for LFU and MFU algorithms.
	Frequency int
}

// Bytes converts Response data structure into bytes array.
func (r Response) Bytes() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(&r)

	return b.Bytes()
}

func (r Response) String() string {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(&r)

	return string(b.Bytes())
}
