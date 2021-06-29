package main

import "math/rand"

const (
	BoardWidth  int = 10
	BoardHeight int = 10
)

type Cell int

const (
	CellDead Cell = iota
	CellLive
)

type Board struct {
	arr [BoardWidth * BoardHeight]Cell
}

func NewBoard() Board {
	board := Board{}
	random := rand.New(rand.NewSource(0))
	for i := range board.arr {
		board.arr[i] = Cell(random.Int63() & 1)
	}
	return board
}

func (b *Board) Cell(x int, y int) Cell {
	return b.arr[y*BoardWidth+x]
}

func getCellWrapAround(board *[BoardWidth * BoardHeight]Cell, index int) int {
	if index < 0 {
		return int(board[BoardWidth*BoardHeight+index])
	}
	return int(board[index%(BoardWidth*BoardHeight)])
}

func (b *Board) Tick() {
	tmp := b.arr
	for i := range b.arr {
		sum := getCellWrapAround(&b.arr, i-BoardWidth-1) +
			getCellWrapAround(&b.arr, i-BoardWidth) +
			getCellWrapAround(&b.arr, i-BoardWidth+1) +
			getCellWrapAround(&b.arr, i-1) +
			getCellWrapAround(&b.arr, i+1) +
			getCellWrapAround(&b.arr, i+BoardWidth-1) +
			getCellWrapAround(&b.arr, i+BoardWidth) +
			getCellWrapAround(&b.arr, i+BoardWidth+1)

		if sum < 2 || sum > 3 {
			tmp[i] = CellDead
		} else if b.arr[i] == CellDead && sum == 3 {
			tmp[i] = CellLive
		}
	}
	b.arr = tmp
}
