package main

import (
	"fmt"
	"time"

	"github.com/robstrong/nest-bt/btpoll"
	"github.com/robstrong/nest-bt/nest"
)

func main() {
	n := nest.NewNestHandler()

	//start BT listener
	bt := btpoll.New(5 * time.Second)
	bt.AddDevice(
		"90:00:DB:3B:9C:C4",
		n.Found,
		n.NotFound,
	)
	err := bt.Start()

	if err != nil {
		fmt.Printf("err starting: %s\n", err)
	}
}
