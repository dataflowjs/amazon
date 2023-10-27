package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

type Instance struct {
	Launcher *launcher.Launcher
	Browser  *rod.Browser
	Page     *rod.Page
	Url      string
	Port     string
}

func (i *Instance) SetupBrowser() {
	log.Println("connecting...")

	i.Port = "10010"

	profileId := "e638a91c-8cac-4644-826b-ce00d8a25364"
	rq := "http://127.0.0.1:10010/api/v1/profile/start?automation=true&puppeteer=true&profileId=" + profileId

	rsp, err := http.Get(rq)
	if err != nil {
		log.Println(err)
	}
	defer rsp.Body.Close()

	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err)
	}

	var mla MLA

	err = json.Unmarshal(b, &mla)
	if err != nil {
		log.Println(err)
	}

	log.Println(mla.Status, mla.Value)

	link, err := url.Parse(mla.Value)
	if err != nil {
		log.Println(err)
	}

	log.Println(link.Port())

	i.Port = link.Port()

	utils.Sleep(2)

	ln, err := launcher.ResolveURL(mla.Value)
	if err != nil {
		log.Println("is not a browser")
		log.Println(err)
	}

	log.Println(ln, err)

	err = rod.Try(func() {
		i.Browser = rod.New().
			NoDefaultDevice().
			ControlURL(ln).
			MustConnect()
	})
	if err != nil {
		log.Println(err)
	}

	i.Browser.Trace(true)

	utils.Sleep(1)

	pages := i.Browser.MustPages()

	i.Page = pages[0]

	go i.Browser.MustHandleAuth("c5l5y", "4ko1f9pp")()
}

func (i *Instance) Process() {
	log.Println("process")
	for _, asin := range ASINs {
		err := i.Page.Navigate("https://sellercentral.amazon.com/product-search?ref=xx_catadd_dnav_xx")
		if err != nil {
			log.Println(err)
		}

		utils.Sleep(2)

		log.Println("search")
		log.Println(i.Page.MustInfo().Title)

		click(i, ".search-input-group", 100, 24)

		utils.Sleep(1)

		i.Page.Mouse.MustClick(proto.InputMouseButtonLeft)

		utils.Sleep(2)

		for _, a := range asin {
			key := getKey(string(a))
			i.Page.Keyboard.Type(key)

			utils.Sleep(0.3)
		}

		utils.Sleep(2)

		click(i, ".search-button", 30, 24)

		utils.Sleep(1)

		i.Page.Mouse.MustClick(proto.InputMouseButtonLeft)

		utils.Sleep(7)

		click(i, ".listing-actions-dropdown", 50, 12)

		utils.Sleep(2)

		click(i, ".listing-actions-dropdown", 50, 46)

		el := i.Page.Timeout(15 * time.Second).MustSearch("a.copy-kat-button")
		text := el.MustText()

		if text == "Sell this product" {
			err = setSheetData(client, spreadsheetID, sheetName, 1, 8, "Product was Pre-Ungated")
			if err != nil {
				log.Println(err)
			}

			err = setCellBackgroundColor(client, spreadsheetID, sheetName, 1, 8, preUngatedColor)
			if err != nil {
				log.Println(err)
			}
		} else if text == "Apply to sell" {
			click(i, ".copy-kat-button", 50, 18)
			utils.Sleep(5)

			pages := i.Browser.MustPages()

			if len(pages) <= 1 {
				return
			}

			page := pages[1]

			page.Race().Element(".a-button-input").MustHandle(func(e *rod.Element) {
				click(i, ".a-button-input", 50, 15)

				utils.Sleep(5)

			}).ElementFunc(func(p *rod.Page) (*rod.Element, error) {
				el := page.MustElement("body")
				s := el.MustText()

				switch true {
				case strings.Contains(s, "Your selling application is approved"):
					err = setSheetData(client, spreadsheetID, sheetName, 1, 8, "Instant Automatic Ungating")
					if err != nil {
						log.Println(err)
					}

					err = setCellBackgroundColor(client, spreadsheetID, sheetName, 1, 8, automaticUngatingColor)
					if err != nil {
						log.Println(err)
					}

				case strings.Contains(s, "Tell us about your products and business"):
					err = setSheetData(client, spreadsheetID, sheetName, 1, 8, "Ungating After Completing Additional Test")
					if err != nil {
						log.Println(err)
					}

					err = setCellBackgroundColor(client, spreadsheetID, sheetName, 1, 8, additionalTestColor)
					if err != nil {
						log.Println(err)
					}

				case strings.Contains(s, "You are requesting approval to sell"):
					err = setSheetData(client, spreadsheetID, sheetName, 1, 8, "Additional Documents Required")
					if err != nil {
						log.Println(err)
					}

					err = setCellBackgroundColor(client, spreadsheetID, sheetName, 1, 8, documentRequiredColor)
					if err != nil {
						log.Println(err)
					}

				}

				return el, errors.New("state not found")
			})
		}
	}
}
