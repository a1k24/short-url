package pkg

import (
	"crypto/md5"
	"fmt"
	"io"
)

const (
	base         int64 = 62
	characterSet       = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Takes an input as string and returns the its base32 md5 hash
func CreateMd5hash(input string) string {
	h := md5.New()
	io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Takes an unsigned integer input and returns a base62 string as hash
func ToBase62(num int64) string {
	encoded := ""
	for num > 0 {
		r := num % base
		num /= base
		encoded = string(characterSet[r]) + encoded

	}
	return encoded
}
