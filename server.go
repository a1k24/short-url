package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type UrlInfo struct {
	LongUrl   string `json:"long_url""`
	UrlId     string `json:"url_id"`
	LinkHash  string `json:"link_hash"`
	Timestamp int64  `json:"timestamp"`
	ShortUrl  string `json:"short_url"`
	UrlMd5    string `json:"url_md5"`
}

type UrlAlreadyExistsError string

func (e UrlAlreadyExistsError) Error() string {
	return fmt.Sprintf("URL with name: %s already exists!", string(e))
}

var Urls []UrlInfo

var counter uint64 = 0
var baseUrl = "localhost:10000"

type ShortUrlRequest struct {
	LongUrl      string `json:"long_url"`
	ForceNewHash bool   `json:"force_new_hash"`
	CustomName   string `json:"custom_name"`
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{id}", redirectHandler)
	router.HandleFunc("/api/shorten", saveHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(baseUrl, router))
}

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var shortUrlRequest ShortUrlRequest
	json.Unmarshal(reqBody, &shortUrlRequest)

	// validate long url
	if !isUrl(shortUrlRequest.LongUrl) {
		http.Error(writer, "Long url is invalid: "+shortUrlRequest.LongUrl, http.StatusBadRequest)
		return
	}

	urlInfo, err := getOrCreateUrlInfo(&shortUrlRequest)
	if nil != err {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	Urls = append(Urls, *urlInfo) // replace this with save in mongo

	json.NewEncoder(writer).Encode(urlInfo)
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func getOrCreateUrlInfo(shortUrlRequest *ShortUrlRequest) (*UrlInfo, error) {
	longUrl := shortUrlRequest.LongUrl
	if "" != shortUrlRequest.CustomName {
		return createCustomUrl(shortUrlRequest)
	}

	var urlInfo *UrlInfo = nil
	if !shortUrlRequest.ForceNewHash {
		urlInfo = findExistingUrlInfo(longUrl)
	}
	if nil == urlInfo {
		urlInfo = createUrlInfo(shortUrlRequest)
	}
	return urlInfo, nil
}

func createCustomUrl(shortUrlRequest *ShortUrlRequest) (*UrlInfo, error) {
	customName := shortUrlRequest.CustomName
	urlInfo := findUrlInfo(customName)
	if nil != urlInfo {
		return nil, UrlAlreadyExistsError(customName)
	}
	return createUrlInfo(shortUrlRequest), nil
}

func findExistingUrlInfo(longUrl string) *UrlInfo {
	// find from mongo
	md5hash := createMd5hash(longUrl)
	for _, urlInfo := range Urls {
		if md5hash == urlInfo.UrlMd5 {
			return &urlInfo
		}
	}
	return nil
}

func createMd5hash(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func createUrlInfo(shortUrlRequest *ShortUrlRequest) *UrlInfo {
	counter++ // assumed threadsafe
	var linkHash = shortUrlRequest.CustomName
	if "" == linkHash {
		linkHash = toBase62(counter)
	}
	return &UrlInfo{
		LongUrl:   shortUrlRequest.LongUrl,
		UrlId:     strconv.FormatUint(counter, 10),
		LinkHash:  linkHash,
		Timestamp: makeTimestamp(),
		ShortUrl:  baseUrl + "/" + linkHash,
		UrlMd5:    createMd5hash(shortUrlRequest.LongUrl),
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func findUrlInfo(urlHash string) *UrlInfo {
	// find from mongo
	for _, urlInfo := range Urls {
		if urlHash == urlInfo.LinkHash {
			return &urlInfo
		}
	}
	return nil
}

func redirectHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Request received.", request)
	vars := mux.Vars(request)
	key := vars["id"]
	urlInfo := findUrlInfo(key)
	if nil == urlInfo {
		http.NotFound(writer, request)
	}
	http.Redirect(writer, request, urlInfo.LongUrl, http.StatusFound)
}

const (
	base         uint64 = 62
	characterSet        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func toBase62(num uint64) string {
	encoded := ""
	for num > 0 {
		r := num % base
		num /= base
		encoded = string(characterSet[r]) + encoded

	}
	return encoded
}

func main() {
	fmt.Println("Server started.")
	handleRequests()
}
