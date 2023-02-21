package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// Algo hash algorithm
type Algo string

const (
	AlgoMD5    Algo = "md5"
	AlgoSha1   Algo = "sha1"
	AlgoSha224 Algo = "sha224"
	AlgoSha256 Algo = "sha256"
	AlgoSha384 Algo = "sha384"
	AlgoSha512 Algo = "sha512"
)

// Hash generates a hash value, expects: MD5, SHA1, SHA224, SHA256, SHA384, SHA512.
func Hash(algo Algo, s string) string {
	var h hash.Hash

	switch algo {
	case AlgoMD5:
		h = md5.New()
	case AlgoSha1:
		h = sha1.New()
	case AlgoSha224:
		h = sha256.New224()
	case AlgoSha256:
		h = sha256.New()
	case AlgoSha384:
		h = sha512.New384()
	case AlgoSha512:
		h = sha512.New()
	default:
		return s
	}

	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

// HMAC generates a keyed hash value, expects: MD5, SHA1, SHA224, SHA256, SHA384, SHA512.
func HMAC(algo Algo, s, key string) string {
	var mac hash.Hash

	switch algo {
	case AlgoMD5:
		mac = hmac.New(md5.New, []byte(key))
	case AlgoSha1:
		mac = hmac.New(sha1.New, []byte(key))
	case AlgoSha224:
		mac = hmac.New(sha256.New224, []byte(key))
	case AlgoSha256:
		mac = hmac.New(sha256.New, []byte(key))
	case AlgoSha384:
		mac = hmac.New(sha512.New384, []byte(key))
	case AlgoSha512:
		mac = hmac.New(sha512.New, []byte(key))
	default:
		return s
	}

	mac.Write([]byte(s))

	return hex.EncodeToString(mac.Sum(nil))
}
