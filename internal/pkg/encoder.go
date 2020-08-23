package pkg

import (
	"crypto/md5"
	"fmt"
	"io"
)

const (
	base         uint64 = 62
	characterSet        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func CreateMd5hash(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ToBase62(num uint64) string {
	encoded := ""
	for num > 0 {
		r := num % base
		num /= base
		encoded = string(characterSet[r]) + encoded

	}
	return encoded
}
