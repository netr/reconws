package main

import (
	"fmt"
	"github.com/netr/reconws"
)

func main() {
	chDone := make(chan bool, 2)
	chRead := make(chan string)
	c := reconws.NewClient().SetChannels(chRead, chDone).OnConnect(onConnect)
	_, err := c.Connect("wss://stream.binance.us:9443/stream?streams=btcusd@kline_15m/btcusd@kline_5m")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case msg := <-chRead:
			fmt.Println(msg)
		case <-chDone:
			return
		}
	}
}

func onConnect() {
	fmt.Println("connected!")
}
