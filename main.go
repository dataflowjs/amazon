package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	file, e := os.OpenFile("file.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0666)
	if e != nil {
		log.Fatalln("Failed to open log file")
	}

	multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multi)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// Replace with the path to your credentials JSON file.
	credentialsFilePath := "credentials.json"

	// Initialize the Google Sheets API client.
	client, err := createGoogleSheetsClient(credentialsFilePath)
	if err != nil {
		log.Fatalf("Unable to create Google Sheets client: %v", err)
	}

	// Replace with the ID of your Google Sheet.
	spreadsheetID := "1OCtJlR3yZRaQwzX7joW0GnwiYUcme84WyMr0w3I0KNY"

	// Replace with the name or index of the sheet you want to access.
	sheetName := "Ungating Sheet Template"

	// Retrieve data from the specified sheet.
	data, err := getSheetData(client, spreadsheetID, sheetName)
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	log.Println("assins")
	for i, row := range data {
		if row[1] == "MANUAL INPUT BY OUR TEAM" || row[1] == "ASIN" {
			continue
		}

		fmt.Println(row[1])

		if i == 10 {
			break
		}
	}

	instance := Instance{}
	instance.SetupBrowser()
}
