package main

import (
	"log"

	"net/http"

	"github.com/z2665/chanmonitor/pkg/chanhttp"
	"github.com/z2665/chanmonitor/pkg/chanmonitor"
)

func main() {
	//if chan le/cap > 50% then it's overflow
	mon := chanmonitor.NewChanMonitor(50)
	c := make(chan int, 10)
	c1 := make(chan int, 10)
	mon.AddChan("example-1", c)
	mon.AddChan("example-2", c1)
	for i := 0; i != 7; i++ {
		c <- i
	}
	chanhttp.NewChanHTTP(nil, mon)
	log.Println(http.ListenAndServe(":6060", nil))
}
