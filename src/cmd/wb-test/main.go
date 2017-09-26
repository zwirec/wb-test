package main

import (
	"log"
	"github.com/zwirec/wb-test/src/cmd/wb-test/counter"
)

func init() {
	log.SetFlags(0)
}

func main() {
	counter.Count()
}
