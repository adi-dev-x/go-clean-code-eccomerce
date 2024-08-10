package main

import (
	"fmt"
)

func main() {
	// Create and initialize a map
	myMap := map[string]int{
		"apple":  1,
		"banana": 2,
		"cherry": 3,
	}

	// Iterate over the map and print the key-value pairs
	fmt.Println("Map iteration:")
	for key, value := range myMap {
		fmt.Printf("%s: %d\n", key, value)
	}
}
