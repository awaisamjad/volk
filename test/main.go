package main

import (
	"fmt"
	"flag"
)

var (
	config = flag.String("config", "", "Path to configu")
)


func main() {
	flag.Parse()
	fmt.Println(*config)
}

