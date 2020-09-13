# short-url
A simple Go implementation of URL shortener using mongo DB as persistent storage.

### Current Features Supported
* Generate Shortened URL given an input URL.
* Validate input URL for correctness.
* Generate customized hash for an input URL. ( vanity url)
* Option to generate multiple hashes for same input URL.
* Obtain click-count for particular URL using the link hash.

### To start server
* Run `go build github.com/a1k24/short-url/cmd/server`
* Run `./server -username <username> -password <password>`

### Using Docker
* Run `docker build --build-arg USERNAME=<username> --build-arg PASSWORD=<password> --build-arg DOMAIN=<domain> -t server . -f build/package/Dockerfile`
* Run `docker run --publish 127.0.0.1:8080:10000 --name test --rm server`
* Latest stable [image](https://hub.docker.com/r/a1k24/short-url)

### Sample CURL:
Create Short URL:
```
curl -XPOST -H 'Content-Type:application/json' "localhost:10000/api/shorten" \
-d '{"long_url": "https://www.google.com"}'
```
Response:
```
{
  "long_url": "https://www.google.com",
  "url_id": "6",
  "link_hash": "6",
  "timestamp": 1598205267552,
  "short_url": "localhost:10000/6",
  "url_md5": "8ffdefbdec956b595d257f0aaeefd623",
  "click_count": 0
}
```
Create Short URL ( force new hash ):
```
curl -XPOST -H 'Content-Type:application/json' "localhost:10000/api/shorten" \
-d '{"long_url": "https://www.google.com", "force_new_hash" : true}'
```
Response:
```
{
  "long_url": "https://www.google.com",
  "url_id": "8",
  "link_hash": "8",
  "timestamp": 1598205267552,
  "short_url": "localhost:10000/8",
  "url_md5": "8ffdefbdec956b595d257f0aaeefd623",
  "click_count": 0
}
```
Create Short URL with custom name:
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
  "url_md5": "8ffdefbdec956b595d257f0aaeefd623",
  "click_count": 0
}
```
Fetch Short URL:
```
curl -XGET http://localhost:10000/api/hello12
```
Response:
```
{
  "long_url": "https://www.google.com",
  "url_id": "7",
  "link_hash": "hello12",
  "timestamp": 1598205266552,
  "short_url": "localhost:10000/hello12",
  "url_md5": "8ffdefbdec956b595d257f0aaeefd623",
  "click_count": 19
}
```

Sample hosted at:
https://acash.dev ( can use this in place of localhost:10000 for the API base URL )

Test shortened links:
* https://acash.dev/hello12
* https://acash.dev/hello2




