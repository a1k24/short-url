# short-url
A simple Go implementation of URL shortener using mongo DB as persistent storage.

### Current Features Supported
* Generate Shortened URL given an input URL.
* Validate input URL for correctness.
* Generate customized hash for an input URL. ( vanity url)
* Option to generate multiple hashes for same input URL.

### To start server
* Run `go build github.com/a1k24/short-url/cmd/server`
* Run `./server`

### Sample Curl:
```
curl -XPOST -H 'Content-Type:application/json' "localhost:10000/api/shorten" \
-d '{"long_url": "https://www.google.com", "custom_name" : "hello12"}'
```
Response:
```
{
  "long_url": "https://www.google.com",
  "url_id": "7",
  "link_hash": "hello12",
  "timestamp": 1598205266552,
  "short_url": "localhost:10000/hello12",
  "url_md5": "8ffdefbdec956b595d257f0aaeefd623"
}
```
