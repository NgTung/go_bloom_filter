package main

import (
	"runtime"

	"demo/src/demo/pkg/api"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	api.Start()
}