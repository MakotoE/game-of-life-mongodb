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

	stopChan := make(chan bool, 1)
	go func() {
		event := termbox.PollEvent()
		if event.Err != nil {
			panic(event.Err)
		}

		if event.Type == termbox.EventInterrupt || event.Type == termbox.EventKey {
			stopChan <- true
		}
	}()

	board := NewArrayBoard()
	defer board.Close()

	for {
		select {
		case <-stopChan:
			return
		default:
		}

		for x := 0; x < BoardWidth; x++ {
			for y := 0; y < BoardHeight; y++ {
				color := termbox.ColorBlack
				cell, err := board.Cell(x, y)
				if err != nil {
					panic(err)
				}
				if cell == CellLive {
					color = termbox.ColorWhite
				}
				termbox.SetBg(x, y, color)
			}
		}

		if err := termbox.Flush(); err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 50)
		if err := board.Tick(); err != nil {
			panic(err)
		}
	}
}
