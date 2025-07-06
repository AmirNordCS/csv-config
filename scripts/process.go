package main

import (
	"encoding/csv"
	"log"
	"os"
)

func main() {
	// CRITICAL CHANGE: File paths must be relative to the working directory.
	// The workflow runs this script from the `csv-config` directory.
	sourceFile := "data/new-config.csv"
	destinationFile := "../csv-firewall/data/all-config.csv"

	// 1. Open the source file for reading
	src, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalf("failed to open source file %s: %s", sourceFile, err)
	}
	defer src.Close()

	// Create a new CSV reader
	reader := csv.NewReader(src)
	// We assume the source file has a header, which we'll read and discard
	// so it's not appended to the destination file.
	_, err = reader.Read()
	if err != nil {
		// This will fail if the source file is empty. Add a check.
		if err.Error() == "EOF" {
			log.Println("Source file is empty or has only a header. Nothing to append.")
			return // Exit gracefully
		}
		log.Fatalf("failed to read header from source file: %s", err)
	}

	// Read all the remaining records from the source file
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read records from source file: %s", err)
	}

	if len(records) == 0 {
		log.Println("No new records to append.")
		return // Exit gracefully if there's no data after the header
	}

	// 2. Open the destination file in append mode
	dest, err := os.OpenFile(destinationFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to open destination file %s: %s", destinationFile, err)
	}
	defer dest.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(dest)
	defer writer.Flush()

	// Write all the records from the source to the destination
	err = writer.WriteAll(records)
	if err != nil {
		log.Fatalf("failed to write records to destination file: %s", err)
	}

	log.Printf("Successfully appended %d records to %s", len(records), destinationFile)
}
