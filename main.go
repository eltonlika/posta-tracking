package main

import (
	"fmt"
	"os"
	"time"

	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
)

var trackingURL = "https://gjurmo.postashqiptare.al/tracking.aspx"
var timeout = time.Second * 8

func main() {
	trackingNumber := os.Args[2]

	bow := surf.NewBrowser()
	bow.SetAttribute(browser.FollowRedirects, true)
	bow.SetUserAgent(agent.Chrome())
	bow.SetTimeout(timeout)

	err := bow.Open(trackingURL)
	if err != nil {
		panic(err)
	}

	fm, _ := bow.Form("#form1")
	fm.Input("txt_barcode", trackingNumber)
	fm.Input("hBarCodes", trackingNumber)
	bow.Click("")
	if fm.Submit() != nil {
		panic(err)
	}

	fmt.Println(bow.Body())
}
