package app

import (
	"encoding/json"
	"github.com/a1k24/short-url/configs"
	"github.com/a1k24/short-url/internal/pkg"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{id}", redirectHandler)
	router.HandleFunc("/api/shorten", saveHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(configs.BaseUrl, router))
}

var Urls []UrlInfo

var counter uint64 = 0

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var shortUrlRequest ShortUrlRequest
	json.Unmarshal(reqBody, &shortUrlRequest)

	// validate long url
	if !pkg.IsUrl(shortUrlRequest.LongUrl) {
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

func getOrCreateUrlInfo(shortUrlRequest *ShortUrlRequest) (*UrlInfo, error) {
	longUrl := shortUrlRequest.LongUrl
	if "" != shortUrlRequest.CustomName {
		return createCustomUrl(shortUrlRequest)
	}

	var urlInfo *UrlInfo = nil
	if !shortUrlRequest.ForceNewHash {
		urlInfo = findUrlInfoForLongUrl(longUrl)
	}
	if nil == urlInfo {
		urlInfo = createUrlInfo(shortUrlRequest)
	}
	return urlInfo, nil
}

func createCustomUrl(shortUrlRequest *ShortUrlRequest) (*UrlInfo, error) {
	customName := shortUrlRequest.CustomName
	urlInfo := findUrlInfoForHash(customName)
	if nil != urlInfo {
		return nil, pkg.UrlAlreadyExistsError(customName)
	}
	return createUrlInfo(shortUrlRequest), nil
}

func findUrlInfoForLongUrl(longUrl string) *UrlInfo {
	// find from mongo
	md5hash := pkg.CreateMd5hash(longUrl)
	for _, urlInfo := range Urls {
		if md5hash == urlInfo.UrlMd5 {
			return &urlInfo
		}
	}
	return nil
}

func createUrlInfo(shortUrlRequest *ShortUrlRequest) *UrlInfo {
	counter++ // assumed threadsafe
	var linkHash = shortUrlRequest.CustomName
	if "" == linkHash {
		linkHash = pkg.ToBase62(counter)
	}
	return &UrlInfo{
		LongUrl:   shortUrlRequest.LongUrl,
		UrlId:     strconv.FormatUint(counter, 10),
		LinkHash:  linkHash,
		Timestamp: pkg.MakeTimestamp(),
		ShortUrl:  configs.BaseUrl + "/" + linkHash,
		UrlMd5:    pkg.CreateMd5hash(shortUrlRequest.LongUrl),
	}
}

func findUrlInfoForHash(linkHash string) *UrlInfo {
	// find from mongo
	for _, urlInfo := range Urls {
		if linkHash == urlInfo.LinkHash {
			return &urlInfo
		}
	}
	return nil
}

func redirectHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]
	urlInfo := findUrlInfoForHash(key)
	if nil == urlInfo {
		http.NotFound(writer, request)
	}
	http.Redirect(writer, request, urlInfo.LongUrl, http.StatusFound)
}
