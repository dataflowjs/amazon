package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

func click(i *Instance, selector string, x, y float64) {
	_, err := i.Page.Timeout(15 * time.Second).Element(selector)
	if err != nil {
		log.Println(err)
	}

	utils.Sleep(3)

	pos, err := i.Page.Eval(`function (selector) {
			let el = document.querySelector(selector);
			var clientRect = el.getBoundingClientRect();
			return {left: clientRect.left + document.body.scrollLeft,
					top: clientRect.top + document.body.scrollTop};
		}`, selector)
	if err != nil {
		log.Println(err)
	}

	posX := pos.Value.Get("left").Num()
	posY := pos.Value.Get("top").Num()

	log.Println(posX, posY)

	i.Page.Mouse.MoveLinear(proto.Point{
		X: posX + x,
		Y: posY + y,
	}, 50)
}

func prettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	r := string(s)
	fmt.Println(r)
}

func JoinExecutePath(path ...string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	newPath := append([]string{}, exPath)
	newPath = append(newPath, path...)

	return filepath.Join(newPath...)
}

func ApiSend(msg, key, channel string) error {
	rq, err := http.NewRequest("GET", "https://api.telegram.org/bot"+key+"/sendMessage", nil)
	if err != nil {
		return err
	}

	rq.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	rq.Header.Set("Host", "api.telegram.org")

	q := url.Values{}
	q.Set("chat_id", channel)
	q.Set("parse_mode", "HTML")
	q.Set("disable_web_page_preview", "true")
	q.Set("text", msg)

	rq.URL.RawQuery = q.Encode()

	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return err
	}

	bin, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(bin))

	return nil
}
