package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cl-parse/changelog"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: cl-parse <changelog-file>")
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	parser := changelog.NewParser()
	entries, err := parser.Parse(string(content))
	if err != nil {
		log.Fatal(err)
	}

	// Convert to JSON for viewing
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}
