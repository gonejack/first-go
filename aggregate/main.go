package main

import (
	"github.com/davecgh/go-spew/spew"
	"time"
)

type Key struct {
	Date int
	Hour int
	//...
}
type Log struct {
	Key
	Bids   int
	Offers int
}

func (lg *Log) aggregate(log Log) {
	lg.Bids += log.Bids
	lg.Offers += log.Offers
}

func main() {
	var m = make(map[Key]Log)
	var c = make(chan Log)

	time.AfterFunc(time.Second*10, func() { close(c) })

	go func() {
		for range time.Tick(time.Second) {
			c <- Log{Bids: 1}
		}
	}()

	for {
		select {
		case log, ok := <-c:
			if !ok {
				return
			}
			sum, exist := m[log.Key]
			if exist {
				sum.aggregate(log)
			} else {
				m[log.Key] = log
			}
		case <-time.After(time.Second):
			for _, lg := range m {
				spew.Dump(lg)
			}
			m = make(map[Key]Log)
		}
	}
}
