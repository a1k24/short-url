package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/a1k24/short-url/configs"
	"github.com/a1k24/short-url/internal/pkg"
)

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{id}", redirectHandler)
	router.HandleFunc("/api/shorten", saveHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(configs.GetBaseUrl(), router))
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
		ShortUrl:  configs.GetBaseUrl() + "/" + linkHash,
		UrlMd5:    pkg.CreateMd5hash(shortUrlRequest.LongUrl),
	}, nil
}

func redirectHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]
	urlInfo, err := FindUrlByLinkHash(key)
	if nil != err {
		log.Println(err)
		http.Error(writer, "Unknown error occurred.", http.StatusInternalServerError)
	}
	if nil == urlInfo {
		http.NotFound(writer, request)
	}
	http.Redirect(writer, request, urlInfo.LongUrl, http.StatusFound)
}
