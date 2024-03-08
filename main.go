package main

import (
	"bufio"
	"enspired/src"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing input file argument")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Could not open file %s: %v", os.Args[1], err)
		os.Exit(1)
	}
	defer file.Close()

	parser := src.NewRoomParser()
	err = parser.IngestAllFromReader(bufio.NewReader(file))
	if err != nil {
		fmt.Printf("Error processing input: %v", err)
		return
	}

	fmt.Printf("%s\n", parser)
}
