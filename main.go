package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/sheets/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

var ASINs []string
var client *sheets.Service
var spreadsheetID = "1OCtJlR3yZRaQwzX7joW0GnwiYUcme84WyMr0w3I0KNY"
var sheetName = "Ungating Sheet Template"
var Db *gorm.DB

func main() {
	credentialsFilePath := "credentials.json"

	var err error
	client, err = createGoogleSheetsClient(credentialsFilePath)
	if err != nil {
		log.Fatalf("Unable to create Google Sheets client: %v", err)
	}

	data, err := getSheetData(client, spreadsheetID, sheetName)
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	log.Println("assins")
	for i, row := range data {
		if i == 0 || i == 1 {
			continue
		}

		ASINs = append(ASINs, row[1].(string))
	}

	instance := Instance{}
	instance.SetupBrowser()
	instance.Process()

	Db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = Db.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}

	err = Db.First(&Stg).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err = Db.Create(&Stg).Error; err != nil {
			log.Println(err)
		}
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/assets", "./assets")

	router.GET("/", IndexHandler)

	router.Run(":8080")
}
