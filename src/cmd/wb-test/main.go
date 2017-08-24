package main

import (
	"bufio"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

// Other funcs, if any.

func main() {

	// Initialization, if any.

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		// Do something with strings here.

	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// Other code, if any.

	total := 0

	// Other code, if any.

	log.Printf("Total: %v", total)
}
