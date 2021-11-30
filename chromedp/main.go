package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	downloadWaitGroup := sync.WaitGroup{}
	downloadWaitGroup.Add(1)
	downloadURL := "https://avatars.githubusercontent.com/u/33149672"

	var requestId network.RequestID

	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			fmt.Printf("EventRequestWillBeSent: %v: %v\n", ev.RequestID, ev.Request.URL)
			if ev.Request.URL == downloadURL {
				requestId = ev.RequestID
			}
		case *network.EventLoadingFinished:
			fmt.Printf("EventLoadingFinished: %v\n", ev.RequestID)
			if ev.RequestID == requestId {
				downloadWaitGroup.Done()
			}
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.Navigate(downloadURL),
	); err != nil {
		log.Fatal(err)
	}

	downloadWaitGroup.Wait()

	var downloadBytes []byte
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		downloadBytes, err = network.GetResponseBody(requestId).Do(ctx)
		return err
	})); err != nil {
		log.Fatal(err)
	}

	downloadDest := fmt.Sprintf("%v/download.png", os.TempDir())
	if err := ioutil.WriteFile(downloadDest, downloadBytes, 0777); err != nil {
		log.Fatal(err)
	}

	log.Printf("Download Complete: %v", downloadDest)
}

// print a specific pdf page.
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
