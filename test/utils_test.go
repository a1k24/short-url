package test

import (
	"github.com/a1k24/short-url/internal/pkg"
	"testing"
	"time"
)

func TestIsUrl(t *testing.T) {
	tests := []struct {
		inputUrl string
		response bool
	}{
		{"https://www.google.com", true},
		{"www.google.com", false},
		{"http://abcde", false},
		{"", false},
	}

	for _, test := range tests {
		if g, w := pkg.IsValidUrl(test.inputUrl), test.response; g != w {
			t.Errorf("%s: is_url = %v, want %v", test.inputUrl, g, w)
		}
	}
}

func TestMakeTimestamp(t *testing.T) {
	timestamp1 := pkg.MakeTimestamp()
	time.Sleep(2 * time.Millisecond)
	timestamp2 := pkg.MakeTimestamp()
	if timestamp1 >= timestamp2 {
		t.Errorf("Got non-incrasing timestamp. T1: %d, T2: %d", timestamp1, timestamp2)
	}
}
