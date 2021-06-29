package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetChar(0, 0, 'a')
	if err := termbox.Flush(); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 2)
}
