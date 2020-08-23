package app

type UrlInfo struct {
	LongUrl   string `json:"long_url""`
	UrlId     string `json:"url_id"`
	LinkHash  string `json:"link_hash"`
	Timestamp int64  `json:"timestamp"`
	ShortUrl  string `json:"short_url"`
	UrlMd5    string `json:"url_md5"`
}

type ShortUrlRequest struct {
	LongUrl      string `json:"long_url"`
	ForceNewHash bool   `json:"force_new_hash"`
	CustomName   string `json:"custom_name"`
}
