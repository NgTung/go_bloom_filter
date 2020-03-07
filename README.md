# Bloom Filter in Golang
A simple bloom filter implementation in Golang run on [Gin framework](https://github.com/gin-gonic/gin).

Original blog (Vietnamese): https://medium.com/@ngtung/bloom-filter-trong-golang-211d00443721

Size of items in demo is 50,000 product item name.
 
#### Source clone:
```bash
$ git clone https://github.com/NgTung/go_bloom_filter.git
```
#### Run app:
```bash
$ go get -u github.com/gin-gonic/gin
$ go get -u github.com/spaolacci/murmur3
$ cd src/demo/cmd/; go run main.go
```
* Query item:
```bash
$ curl -X GET -H "Accept: application/json" http://localhost:8080/product?item=rinne
```
