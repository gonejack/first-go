package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"os/signal"
	"time"
)

func main() {
	session, sessionEvents, err := zk.Connect([]string{"192.168.10.162:2181", "192.168.10.163:2181", "192.168.10.164:2181"}, time.Second)

	spew.Dump(sessionEvents, err)

	c := make(chan os.Signal, 1)

	signal.Notify(c)

	_, stat, childEvents, err := session.ChildrenW("/first")

	spew.Dump(stat, err)

	for {
		select {
		case ev, ok := <-childEvents:
			if ok {
				fmt.Printf("%s %v\n", ev, ok)
			} else {
				_, stat, childEvents, err = session.ChildrenW("/first")
			}
		case <-time.After(time.Second * 5):
			_, stat, err := session.Get("/first/data")
			spew.Dump(stat, err)
			spew.Dump(session.Set("/first/data", []byte(fmt.Sprintf("%s", time.Now())), stat.Version))
			//spew.Dump(session.Create("/first/dd", nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll)))
		}

		time.Sleep(time.Second)
	}
}
