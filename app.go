package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

type Instance struct {
	Launcher   *launcher.Launcher
	Browser    *rod.Browser
	Page       *rod.Page
	UploadPage *rod.Page
	Url        string
	Port       string
}

func (i *Instance) SetupBrowser() {
	// var err error
	// var ln string
	log.Println("connecting...")

	i.Port = "10010"

	// for idx := 10008; idx < 10011; idx++ {
	// 	// i.Port = strconv.Itoa(idx)
	// 	log.Println(i.Port)

	// 	conn, err := net.DialTimeout("tcp", net.JoinHostPort("", i.Port), 2*time.Millisecond)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	conn.Close()

	// 	log.Println("connecting...")

	// 	ln, err = launcher.ResolveURL(i.Port)
	// 	if err != nil {
	// 		log.Println("is not a browser")
	// 		time.Sleep(5 * time.Millisecond)
	// 		continue
	// 	}

	// 	log.Println(ln)
	// 	log.Println("connected to port: " + i.Port)
	// 	err = rod.Try(func() {
	// 		i.Browser = rod.New().
	// 			NoDefaultDevice().
	// 			ControlURL(ln).
	// 			MustConnect()
	// 	})
	// 	if err == nil {
	// 		break
	// 	}
	// }

	profileId := "e638a91c-8cac-4644-826b-ce00d8a25364"
	rq := "http://127.0.0.1:10010/api/v1/profile/start?automation=true&profileId=" + profileId

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

	utils.Sleep(5)

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

	for _, p := range pages {
		log.Println(p.MustInfo().Title)
	}

	// i.Page =
}

func (i *Instance) Process() {
}
