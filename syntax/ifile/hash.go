package ifile

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
)

// 对图片有溢出问题
func FileHashSha256(reader io.Reader) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		log.Println(err)
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func FileHashMd5(reader io.Reader) string {
	m := md5.New()
	_, err := io.Copy(m, reader)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(m.Sum(nil))
}

// #SHA1
func FileHashSHA1(reader io.Reader) string {
	h := sha1.New()
	_, err := io.Copy(h, reader)
	// log.Println(h, "hg",h.Size())
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
