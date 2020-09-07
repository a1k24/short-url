package test

import (
	"os"
	"testing"

	"github.com/a1k24/short-url/internal/app"
	"github.com/a1k24/short-url/internal/pkg"
)

var info = app.UrlInfo{
	LongUrl:   "https://www.google.com",
	UrlId:     "dummy",
	LinkHash:  "dummy",
	Timestamp: 0,
	ShortUrl:  "dummy",
	UrlMd5:    "dummy",
}

func TestDao(t *testing.T) {
	_, cancel := pkg.CreateConnection(os.Getenv("MONGO_URL")) // ensure mongo client is created at start
	defer cancel()

	findByLinkHashAndAssert(t, false)

	findByMd5AndAssert(t, false)

	saveAndAssert(t)

	findByLinkHashAndAssert(t, true)

	findByMd5AndAssert(t, true)

	removeByLinkHashAndAssert(t)
}

func removeByLinkHashAndAssert(t *testing.T) {
	deleteResult, err := app.RemoveUrlFromDB(info.LinkHash)
	if nil != err {
		t.Error("Failed to remove from DB", err)
	}
	if g, w := deleteResult.DeletedCount, int64(1); g != w {
		t.Errorf("linkhash: %s, delete_count = %v, want %v", info.LinkHash, g, w)
	}
}

func findByMd5AndAssert(t *testing.T, present bool) {
	urlInfo, err := app.FindUrlByUrlMd5(info.UrlMd5)
	if nil != err {
		t.Error("Failed to find by Md5", err)
	}
	if !present {
		if urlInfo != nil {
			t.Error("Found random data by Md5", err)
		}
		return
	}
	if g, w := *urlInfo, info; g != w {
		t.Errorf("md5: %s, urlinfo = %v, want %v", info.UrlMd5, g, w)
	}
}

func incrementClickCount(t *testing.T, present bool) {
	app.IncrementClickCount(info.LinkHash)
	urlInfo, err := app.FindUrlByLinkHash(info.LinkHash)
	if nil != err {
		t.Error("Failed to find by LinkHash", err)
	}
	if !present {
		if urlInfo != nil {
			t.Error("Found random data by Md5", err)
		}
		return
	}
	if g, w := urlInfo.ClickCount, 1; g != w {
		t.Errorf("url_info: %s, count = %v, want %v", info.LinkHash, g, w)
	}
}

func findByLinkHashAndAssert(t *testing.T, present bool) {
	urlInfo, err := app.FindUrlByLinkHash(info.LinkHash)
	if nil != err {
		t.Error("Failed to find by linkHash", err)
	}
	if !present {
		if urlInfo != nil {
			t.Error("Found random data by linkHash", err)
		}
		return
	}
	if g, w := *urlInfo, info; g != w {
		t.Errorf("linkhash: %s, urlinfo = %v, want %v", info.LinkHash, g, w)
	}
}

func saveAndAssert(t *testing.T) {
	result, err := app.SaveUrlToDB(&info)
	if nil != err {
		t.Error("Failed to save to DB", err)
	}
	if nil == result || "" == result.InsertedID {
		t.Error("Failed to save to DB. Invalid inserted ID", err)
	}
}
