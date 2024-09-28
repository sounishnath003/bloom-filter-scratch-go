package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	Store    []bool            // to Store Bitset
	Size     int32             // Size of the filters
	mHashers []murmur3.Hash128 // HashFunctions to improve the probabilistic accuracy
}

// NewBloomFilter helps to create a new bloom.Filter with the given size.
//
// It initializes the bloom.Store with the size and generates the
// murmur3 hash functions with different seeds.
//
// The number of hash functions is fixed to 8.
//
// Returns a new bloom.Filter.
func NewBloomFilter(size int32) *BloomFilter {
	// Generating the Seed and MHashers (Murmur 128 style)
	//
	return &BloomFilter{
		Store: make([]bool, size),
		Size:  size,
		mHashers: []murmur3.Hash128{
			murmur3.New128WithSeed(uint32(11)),
			murmur3.New128WithSeed(uint32(31)),
			murmur3.New128WithSeed(uint32(131)),
			murmur3.New128WithSeed(uint32(989)),
			murmur3.New128WithSeed(uint32(1919)),
			murmur3.New128WithSeed(uint32(2007)),
			murmur3.New128WithSeed(uint32(31313)),
			murmur3.New128WithSeed(uint32(9281917)),
		},
	}
}

// Info returns the required information for the bloom.Filter configuration
func (bf *BloomFilter) Info() map[string]any {
	return map[string]any{
		"size":           bf.Size,
		"totalHashFuncs": len(bf.mHashers),
	}
}

// ComputeMurmurHash computes and returns the murmur has of a query string `key`
// you have to module it with the bloom.Size to set the index True in bloom.Store
//
// Non-cryptographic hash, fast and efficient and implementation specific.
func (bf *BloomFilter) ComputeMurmurHash(key string, hashFn int) uint64 {
	bf.mHashers[hashFn].Write([]byte(key))
	val, _ := bf.mHashers[hashFn].Sum128()
	bf.mHashers[hashFn].Reset()
	return val
}

// Add helps to add add given key into the bloom.Store. Remember it does not store the
// actual keys. rather it a probabilistic representtal of their presence.
func (bf *BloomFilter) Add(key string) {
	// index := bf.ComputeMurmurHash(key) % uint64(bf.Size)
	// bf.Store[index] = true
	// Utilizing all has functions
	for i := 0; i < len(bf.mHashers); i++ {
		index := bf.ComputeMurmurHash(key, i) % uint64(bf.Size)
		bf.Store[index] = true
	}
}

// Exists helps to lookup if the key present in the bloom.Store.
// In real, the key might not be present even if the return is true. as it
// works as a probabilistic estimation of finding the presence.
func (bf *BloomFilter) Exists(key string) (uint64, bool) {
	// index := bf.ComputeMurmurHash(key) % uint64(bf.Size)
	// return index, bf.Store[index]

	for i := 0; i < len(bf.mHashers); i++ {
		index := bf.ComputeMurmurHash(key, i) % uint64(bf.Size)
		if !bf.Store[index] {
			return index, false
		}
	}
	return 0, true
}

// define global contexts
var wg sync.WaitGroup
var testResultsOutChan chan map[string]any

func main() {
	BloomFilterSize := 100_000
	testResultsOutChan = make(chan map[string]any, BloomFilterSize)

	// Generate a dataset.
	dataset, trainDataset, testDataset := generateDataset(20_000) // 20K
	log.Println("total dataset size:", len(dataset))
	log.Println("total train.dataset size:", len(trainDataset))
	log.Println("total test.dataset size:", len(testDataset))
	log.Println("invoking the test.cases into goroutines...")

	// Dynamically change the bloomFilter size
	for bfsize := 1000; bfsize <= BloomFilterSize; bfsize += 10000 {
		wg.Add(1)
		go PerformTests(bfsize, dataset, trainDataset, testDataset)
	}

	wg.Wait()
	close(testResultsOutChan)

	var testResults []map[string]any

	for testResult := range testResultsOutChan {
		bytes, _ := json.Marshal(testResult)
		fmt.Println(string(bytes))
		testResults = append(testResults, testResult)
	}

	// Find the minimum percentageFalsePositives and print
	var bestBloomFilter map[string]any
	for _, testResult := range testResults {
		if bestBloomFilter == nil {
			bestBloomFilter = testResult
		} else {
			if testResult["percentageFalsePositives"].(float64) < bestBloomFilter["percentageFalsePositives"].(float64) {
				bestBloomFilter = testResult
			}
		}
	}

	// Print the best bloom filter.
	bytes, _ := json.Marshal(bestBloomFilter)
	fmt.Println("========= BEST BLOOM FILTER =========")
	fmt.Println(string(bytes))
	fmt.Println("========= BEST BLOOM FILTER =========")
}

func PerformTests(bloomSize int, dataset []string, trainDataset map[string]bool, testDataset map[string]bool) {
	defer wg.Done()
	// Define a bloom.
	bloom := NewBloomFilter(int32(bloomSize))

	// Add the keys into bloom.
	for key := range trainDataset {
		bloom.Add(key)
	}

	// Declare the true positives and falsePositives.
	truePositives := 0
	falsePositives := 0

	// traverse the dataset from all dataset
	for _, key := range dataset {
		_, bloomExists := bloom.Exists(key)

		// If bloom says YES
		if bloomExists {
			if _, ok := trainDataset[key]; ok {
				truePositives++
			}
			if _, ok := testDataset[key]; ok {
				falsePositives++
			}
		}
	}

	testResult := map[string]any{
		"bloomInfo":                bloom.Info(),
		"falsePositives":           falsePositives,
		"percentageFalsePositives": 100 * float64(falsePositives) / float64(len(dataset)),
	}
	testResultsOutChan <- testResult
}

// generateDataset helps to generate a list of random strings.
//
// Returns dataset, trainDatasetMap, testDatasetMap
func generateDataset(n int) (
	[]string, map[string]bool, map[string]bool,
) {
	// Dataset.
	dataset := make([]string, 1)
	// Split the dataset.
	trainDataset := make(map[string]bool)
	testDataset := make(map[string]bool)

	for i := 0; i < int(n/2); i++ {
		id := uuid.New()
		dataset = append(dataset, id.String())
		trainDataset[id.String()] = true
	}

	for i := int(n/2) + 1; i < n; i++ {
		id := uuid.New()
		dataset = append(dataset, id.String())
		testDataset[id.String()] = true
	}

	return dataset, trainDataset, testDataset
}
