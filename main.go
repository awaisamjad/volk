package main

import (
	"log"
)

func init() {
	// Configure logger to include file name and line number
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}
func main() {
}
