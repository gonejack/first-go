package race

import (
	"fmt"
	"time"
)

func Run() {
	i := 0

	go func() {
		for {
			i++ // write i

			time.Sleep(time.Nanosecond)
		}
	}()

	for {
		time.Sleep(time.Second)
		fmt.Println(i) // read i
	}
}
