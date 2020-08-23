package pkg

import (
	"net/http"
	"net/url"
	"time"
)

var validStatusCodes = map[int]struct{}{http.StatusOK: {}, http.StatusAccepted: {}, http.StatusMovedPermanently: {}}

// Takes inputUrl as string and validates for schema, host and if URL returns proper response on HEAD
func IsValidUrl(inputUrl string) bool {
	u, err := url.Parse(inputUrl)
	if nil != err || "" == u.Scheme || "" == u.Host {
		return false
	}
	head, err := http.Head(inputUrl)
	if nil != err {
		return false
	}
	_, ok := validStatusCodes[head.StatusCode]
	return ok
}

// Returns current time in milliseconds
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
