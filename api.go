package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func createGoogleSheetsClient(credentialsFilePath string) (*sheets.Service, error) {
	ctx := context.Background()

	creds, err := google.CredentialsFromJSON(ctx, []byte(ReadFile(credentialsFilePath)), sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	client, err := sheets.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getSheetData(client *sheets.Service, spreadsheetID, sheetName string) ([][]interface{}, error) {
	readRange := sheetName

	resp, err := client.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

func ReadFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}
	return string(data)
}

func setSheetData(client *sheets.Service, spreadsheetID, sheetName string, column, row int, value string) error {
	columnLetter := getColumnLetter(column)

	cellRange := fmt.Sprintf("%s%d", columnLetter, row)

	values := [][]interface{}{{value}}
	data := &sheets.ValueRange{
		Values: values,
	}

	_, err := client.Spreadsheets.Values.Update(spreadsheetID, sheetName+"!"+cellRange, data).
		ValueInputOption("RAW").Do()
	return err
}

func getColumnLetter(column int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if column <= 0 || column > 26 {
		return ""
	}
	return string(alphabet[column-1])
}

func setCellBackgroundColor(client *sheets.Service, spreadsheetID, sheetName string, column, row int, color *sheets.Color) error {
	sheet, err := getSheetInfo(client, spreadsheetID, sheetName)
	if err != nil {
		return err
	}

	gridRange := &sheets.GridRange{
		SheetId:          sheet.Properties.SheetId,
		StartRowIndex:    int64(row) - 1,
		EndRowIndex:      int64(row),
		StartColumnIndex: int64(column) - 1,
		EndColumnIndex:   int64(column),
	}

	log.Println(sheet.Properties.SheetId)

	cellData := &sheets.CellData{
		UserEnteredFormat: &sheets.CellFormat{
			BackgroundColor: color,
		},
	}

	rows := []*sheets.RowData{
		{
			Values: []*sheets.CellData{cellData},
		},
	}

	updateRequest := &sheets.UpdateCellsRequest{
		Range:  gridRange,
		Fields: "userEnteredFormat.backgroundColor",
		Rows:   rows,
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				UpdateCells: updateRequest,
			},
		},
	}

	_, err = client.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Do()
	return err
}

func getSheetInfo(client *sheets.Service, spreadsheetID, sheetName string) (*sheets.Sheet, error) {
	spreadsheet, err := client.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return nil, err
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet, nil
		}
	}

	return nil, fmt.Errorf("Sheet '%s' not found", sheetName)
}
