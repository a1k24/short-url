package test

import (
	"strings"
	"testing"

	"github.com/a1k24/short-url/internal/pkg"
)

func TestCreateMd5hash(t *testing.T) {
	md5Alpha := "0123456789abcdef"
	var tests = []struct {
		input    string
		response string
	}{
		{"https://www.youtube.com", "d245406cb6c9f36be9064c92c34e12e1"},
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
	}

	for _, test := range tests {
		hash := pkg.CreateMd5hash(test.input)
		if g, w := hash, test.response; g != w {
			t.Errorf("%s: hash = %s, want %s", test.input, g, w)
		}

		if g, w := len(hash), 32; g != w {
			t.Errorf("%s: hash_len = %d, want %d", test.input, g, w)
		}

		for _, char := range hash {
			if !strings.ContainsRune(md5Alpha, rune(char)) {
				t.Errorf("Invalid character %s in md5 hash", string(char))
			}
		}
	}
}

func TestToBase62(t *testing.T) {
	base62Alpha := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var tests = []struct {
		input    int64
		response string
	}{
		{0, ""},
		{1, "1"},
		{1<<63 - 1, "AzL8n0Y58m7"},
	}

	for _, test := range tests {

		hash := pkg.ToBase62(test.input)
		if g, w := hash, test.response; g != w {
			t.Errorf("%d: hash = %s, want %s", test.input, g, w)
		}

		for _, char := range hash {
			if !strings.ContainsRune(base62Alpha, rune(char)) {
				t.Errorf("Invalid character %s in md5 hash", string(char))
			}
		}
	}
}
