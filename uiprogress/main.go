package main

import (
	"crypto/rand"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)
import "github.com/gosuri/uiprogress"

func main() {
	waitTime := time.Millisecond * 200
	p := uiprogress.New()
	p.Start()

	var wg sync.WaitGroup

	bar1 := p.AddBar(20).AppendCompleted().PrependElapsed()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar1.Incr() {
			time.Sleep(waitTime)
		}
		fmt.Fprintln(p.Bypass(), "Bar1 finished")
	}()

	bar2 := p.AddBar(40).AppendCompleted().PrependElapsed()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar2.Incr() {
			time.Sleep(waitTime)
		}
		fmt.Fprintln(p.Bypass(), "Bar2 finished")
	}()

	time.Sleep(time.Second)
	bar3 := p.AddBar(20).PrependElapsed().AppendCompleted()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar3.Incr() {
			time.Sleep(waitTime)
		}
		fmt.Fprintln(p.Bypass(), "Bar3 finished")
	}()

	wg.Wait()
}

func bar1() {
	var limit int64 = 1024 * 1024 * 500
	// we will copy 200 Mb from /dev/rand to /dev/null
	reader := io.LimitReader(rand.Reader, limit)
	writer := ioutil.Discard

	// start new bar
	bar := pb.Full.Start64(limit)
	bar.Set("prefix", " ")
	bar.SetMaxWidth(150)
	// create proxy reader
	barReader := bar.NewProxyReader(reader)
	// copy from proxy reader
	io.Copy(writer, barReader)
	// finish bar
	bar.Finish()
}

func bar2() {
	req, _ := http.NewRequest("GET", "https://dl.google.com/go/go1.14.2.src.tar.gz", nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(ioutil.Discard, bar), resp.Body)

	req, _ = http.NewRequest("GET", "https://dl.google.com/go/go1.14.2.src.tar.gz", nil)
	resp, _ = http.DefaultClient.Do(req)
	defer resp.Body.Close()

	bar = progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(ioutil.Discard, bar), resp.Body)
}
