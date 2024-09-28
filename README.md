# Bloom Filter Implementation in Golang

This repository provides a Golang implementation of a Bloom Filter, a probabilistic data structure used to determine if an element is a member of a set. The implementation includes.

## Features 
* **Core Bloom Filter Structure**: A basic Bloom Filter structure with configurable parameters like the number of bits and hash functions.

* **Hash Function Implementations**: Multiple hash functions are included, such as MurmurHash3, DJB2, and a custom implementation.

* **Insertion and Lookup Operations**: Efficient methods for adding elements to the filter and checking for membership.

* **False Positive Rate Calculation**: A function to calculate the estimated false positive rate based on the filter's parameters.

* **Testing and Benchmarking**: Unit tests and benchmarks to ensure correctness and performance.


## How to run

* Makefile:
```bash
install:
	go mod tidy
	go mod download
	go mod verify

run:
	go version
	go run bloom.go
```

* Install:
```bash
make install
```

* Dev:
```bash
make run
```

## Output

```bash
go version
go version go1.23.1 darwin/arm64
go run bloom.go
2024/09/28 21:20:57 total dataset size: 20000
2024/09/28 21:20:57 total train.dataset size: 10000
2024/09/28 21:20:57 total test.dataset size: 9999
2024/09/28 21:20:57 invoking the test.cases into goroutines...
{"bloomInfo":{"size":71000,"totalHashFuncs":8},"falsePositives":457,"percentageFalsePositives":2.285}
{"bloomInfo":{"size":61000,"totalHashFuncs":8},"falsePositives":784,"percentageFalsePositives":3.92}
{"bloomInfo":{"size":81000,"totalHashFuncs":8},"falsePositives":254,"percentageFalsePositives":1.27}
{"bloomInfo":{"size":21000,"totalHashFuncs":8},"falsePositives":8344,"percentageFalsePositives":41.72}
{"bloomInfo":{"size":91000,"totalHashFuncs":8},"falsePositives":139,"percentageFalsePositives":0.695}
{"bloomInfo":{"size":1000,"totalHashFuncs":8},"falsePositives":9999,"percentageFalsePositives":49.995}
{"bloomInfo":{"size":31000,"totalHashFuncs":8},"falsePositives":5236,"percentageFalsePositives":26.18}
{"bloomInfo":{"size":11000,"totalHashFuncs":8},"falsePositives":9939,"percentageFalsePositives":49.695}
{"bloomInfo":{"size":41000,"totalHashFuncs":8},"falsePositives":2895,"percentageFalsePositives":14.475}
{"bloomInfo":{"size":51000,"totalHashFuncs":8},"falsePositives":1573,"percentageFalsePositives":7.865}
========= BEST BLOOM FILTER =========
{"bloomInfo":{"size":91000,"totalHashFuncs":8},"falsePositives":139,"percentageFalsePositives":0.695}
========= BEST BLOOM FILTER =========
```