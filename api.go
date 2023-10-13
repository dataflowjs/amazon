package main

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// createGoogleSheetsClient creates a Google Sheets API client using the provided credentials file.
func createGoogleSheetsClient(credentialsFilePath string) (*sheets.Service, error) {
	ctx := context.Background()

	// Read credentials file.
	creds, err := google.CredentialsFromJSON(ctx, []byte(ReadFile(credentialsFilePath)), sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	// Create Google Sheets API client.
	client, err := sheets.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return client, nil
}

// getSheetData retrieves data from a specific sheet in the Google Sheet.
func getSheetData(client *sheets.Service, spreadsheetID, sheetName string) ([][]interface{}, error) {
	// Specify the range to retrieve data from (e.g., "A1:B10").
	readRange := sheetName

	// Retrieve data from the specified range.
	resp, err := client.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

// ReadFile reads the content of a file and returns it as a string.
func ReadFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}
	return string(data)
}
