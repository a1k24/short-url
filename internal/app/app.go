package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/allegro/bigcache"

	"github.com/a1k24/short-url/configs"
	"github.com/a1k24/short-url/internal/pkg"
)

var redirectCache *bigcache.BigCache

func HandleRequests(cache *bigcache.BigCache) {
	redirectCache = cache
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{id}", redirectHandler)
	router.HandleFunc("/api/{id}", fetchHandler).Methods("GET")
	router.HandleFunc("/api/shorten", saveHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(configs.GetBaseUrl(), router))
}

func fetchHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]
	urlInfo, err := FindUrlByLinkHash(key)
	if nil != err {
		log.Println("Failed to fetch UrlInfo: ", key, err)
		http.Error(writer, "Failed to fetch UrlInfo.", http.StatusInternalServerError)
		return
	}
	if nil == urlInfo {
		log.Println("Could not find linkHash: ", key)
		http.Error(writer, "Could not find linkHash", http.StatusNotFound)
		return
	}
	json.NewEncoder(writer).Encode(urlInfo)
}

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var shortUrlRequest ShortUrlRequest
	json.Unmarshal(reqBody, &shortUrlRequest)

	// validate long url
	if !pkg.IsValidUrl(shortUrlRequest.LongUrl) {
		http.Error(writer, "Long url is invalid: "+shortUrlRequest.LongUrl, http.StatusBadRequest)
		return
	}

	if "" != shortUrlRequest.CustomName {
		handleCustomNameRequest(&shortUrlRequest, writer)
	} else {
		handleShortenRequest(&shortUrlRequest, writer)
	}
}

func handleCustomNameRequest(shortUrlRequest *ShortUrlRequest, writer http.ResponseWriter) {
	urlInfo, err := createCustomUrl(shortUrlRequest)
	if nil != err {
		log.Println("Failed to Create custom Url", err)
		http.Error(writer, fmt.Sprintf("Failed to create short Url: %s", err), http.StatusBadRequest)
		return
	}
	_, err = SaveUrlToDB(urlInfo)
	if nil != err {
		log.Println("Failed to save to DB", err)
		http.Error(writer, "Failed to create short Url.", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(urlInfo)
}

func handleShortenRequest(shortUrlRequest *ShortUrlRequest, writer http.ResponseWriter) {
	longUrl := shortUrlRequest.LongUrl
	var urlInfo *UrlInfo = nil
	var err error = nil
	if !shortUrlRequest.ForceNewHash {
		urlInfo = findUrlInfoForLongUrl(longUrl)
	}
	if nil != urlInfo {
		//return existing url response
		json.NewEncoder(writer).Encode(urlInfo)
		return
	}
	urlInfo, err = createUrlInfo(shortUrlRequest)

	if nil != err {
		log.Println("Failed to create short Url.", err)
		http.Error(writer, "Failed to create short Url.", http.StatusBadRequest)
		return
	}

	//persist in DB
	_, err = SaveUrlToDB(urlInfo)
	if err != nil {
		log.Println("Failed to save to DB", err)
		http.Error(writer, "Failed to create short Url.", http.StatusInternalServerError)
	}
	json.NewEncoder(writer).Encode(urlInfo)
}

func createCustomUrl(shortUrlRequest *ShortUrlRequest) (*UrlInfo, error) {
	customName := shortUrlRequest.CustomName
	urlInfo, err := FindUrlByLinkHash(customName)
	if nil != err {
		return nil, err
	}
	if nil != urlInfo {
		return nil, pkg.UrlAlreadyExistsError(customName)
	}
	info, err := createUrlInfo(shortUrlRequest)
	return info, err
}

func findUrlInfoForLongUrl(longUrl string) *UrlInfo {
	md5hash := pkg.CreateMd5hash(longUrl)
	urlInfo, err := FindUrlByUrlMd5(md5hash)
	if nil != err {
		log.Println(err) // silently ignoring error and returning nil, should not break creation flow
		return nil
	}
	return urlInfo
}

func createUrlInfo(shortUrlRequest *ShortUrlRequest) (*UrlInfo, error) {
	sequence, err := GenerateNextSequence("link_sequence")
	if nil != err {
		log.Println(err)
		return nil, err
	}
	var linkHash = shortUrlRequest.CustomName
	if "" == linkHash {
		linkHash = pkg.ToBase62(sequence)
	}
	return &UrlInfo{
		LongUrl:   shortUrlRequest.LongUrl,
		UrlId:     strconv.FormatInt(sequence, 10),
		LinkHash:  linkHash,
		Timestamp: pkg.MakeTimestamp(),
		ShortUrl:  configs.GetDomain() + "/" + linkHash,
		UrlMd5:    pkg.CreateMd5hash(shortUrlRequest.LongUrl),
	}, nil
}

func redirectHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	linkHash := vars["id"]

	if entry, err := redirectCache.Get(linkHash); err == nil {
		incrementCountAndRedirect(writer, request, linkHash, string(entry))
		return
	}

	log.Println("Missed cache. LinkHash: ", linkHash)
	urlInfo, err := FindUrlByLinkHash(linkHash)
	if nil != err {
		log.Println(err)
		http.Error(writer, "Unknown error occurred.", http.StatusInternalServerError)
		return
	}
	if nil == urlInfo {
		http.NotFound(writer, request)
		return
	}
	redirectCache.Set(linkHash, []byte(urlInfo.LongUrl))
	incrementCountAndRedirect(writer, request, linkHash, urlInfo.LongUrl)
}

func incrementCountAndRedirect(writer http.ResponseWriter, request *http.Request, linkHash string, longUrl string) {
	go IncrementClickCount(linkHash)
	http.Redirect(writer, request, longUrl, http.StatusFound)
}
