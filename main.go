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

	board, err := NewDatabaseBoard()
	if err != nil {
		panic(err)
	}
	defer board.Close()

	for {
		select {
		case <-stopChan:
			return
		default:
		}

		cells, err := board.Cells()
		if err != nil {
			panic(err)
		}

		for x := 0; x < BoardWidth; x++ {
			for y := 0; y < BoardHeight; y++ {
				color := termbox.ColorBlack
				if cells.Get(x, y) == CellLive {
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
