package api

import (
	"fmt"

	"hash"
	"hash/fnv"

	"demo/src/demo/util/file"
	"demo/src/demo/util/filters"
	"github.com/gin-gonic/gin"
	"github.com/spaolacci/murmur3"
)

const DataFileName = "product_name.txt"

var bloomFilter filters.BloomFilters

func Start() {
	server := gin.Default()

	bloomFilter = InitBloomFilter()

	server.Use(RequestValidator, BloomFilter)

	productNames := file.LoadFileContent(DataFileName)

	server.GET("/product", func(ctx *gin.Context) {
		item, _ := ctx.GetQuery("item")
		for _, pName := range productNames {
			if pName == item {
				ctx.JSON(200, gin.H{"result": fmt.Sprintf("[%v] founded in file content", item)})
				return
			}
		}
		ctx.JSON(404, gin.H{"result": "item not found"})
	})

	_ = server.Run() // listen and serve on 0.0.0.0:8080
}

// Request validator middleware
func RequestValidator(ctx *gin.Context) {
	if _, ok := ctx.GetQuery("item"); !ok {
		ctx.JSON(400, gin.H{"error": "item name can not be empty",})
		ctx.Abort()
	}
	ctx.Next()
}

// Bloom Filter middleware
func BloomFilter(ctx *gin.Context) {
	item, _ := ctx.GetQuery("item")
	if bloomFilter == nil || bloomFilter.MightContain(item) {
		ctx.Next() // The item possibly in the set, pass to check in the next step
		return
	}
	// The item definitely NOT in the set, return not found
	ctx.JSON(404, gin.H{"result": fmt.Sprintf("[%v] not found, filtered by filter", item)})
	ctx.Abort()
}

func InitBloomFilter() filters.BloomFilters {
	// Hash function should be non-cryptographic
	hashFunctions:= []hash.Hash {
		fnv.New64(),
		fnv.New64a(),
		murmur3.New64(),
	}

	// Setup a bloomFilter filter with n = 50,000 and k = 3 (but k should calculated by following formula: k=m/n*ln2
	bloomFilter = filters.InitBloomFilter(50000, hashFunctions)

	// Load data from file to bloom filter bit set
	loadData2Filter(bloomFilter)

	return bloomFilter
}

func loadData2Filter(bloom filters.BloomFilters) {
	productNames := file.LoadFileContent(DataFileName)
	for _, d := range productNames {
		bloom.Add(d)
	}
}