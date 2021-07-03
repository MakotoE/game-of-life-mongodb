package main

import (
	"math/rand"
	"time"
)

const (
	BoardWidth  int = 10
	BoardHeight int = 10
)

type Cell int

const (
	CellDead Cell = iota
	CellLive
)

type Board interface {
	Close() error
	Cells() (CellArray, error)
	Set(x int, y int, cell Cell) error
	Tick() error
}

type CellArray struct {
	Arr [BoardWidth * BoardHeight]Cell
}

func NewCellArray(setLive [][2]int) CellArray {
	cellArray := CellArray{}
	for _, coordinate := range setLive {
		cellArray.Arr[coordinate[1]*BoardWidth+coordinate[0]] = CellLive
	}
	return cellArray
}

func (c *CellArray) Get(x int, y int) Cell {
	return c.Arr[y*BoardWidth+x]
}

type ArrayBoard struct {
	arr [BoardWidth * BoardHeight]Cell
}

var _ Board = &ArrayBoard{}

func NewArrayBoard() ArrayBoard {
	board := ArrayBoard{}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range board.arr {
		board.arr[i] = Cell(random.Int63() & 1)
	}
	return board
}

func (b *ArrayBoard) Close() error {
	return nil
}

func (b *ArrayBoard) Cells() (CellArray, error) {
	return CellArray{Arr: b.arr}, nil
}

func (b *ArrayBoard) Set(x int, y int, cell Cell) error {
	b.arr[y*BoardWidth+x] = cell
	return nil
}

func wrapAroundIndex(index int) int {
	if index < 0 {
		return BoardWidth*BoardHeight + index
	}
	return index % (BoardWidth * BoardHeight)
}

func (b *ArrayBoard) Tick() error {
	tmp := b.arr
	for i := range b.arr {
		sum := b.arr[wrapAroundIndex(i-BoardWidth-1)] +
			b.arr[wrapAroundIndex(i-BoardWidth)] +
			b.arr[wrapAroundIndex(i-BoardWidth+1)] +
			b.arr[wrapAroundIndex(i-1)] +
			b.arr[wrapAroundIndex(i+1)] +
			b.arr[wrapAroundIndex(i+BoardWidth-1)] +
			b.arr[wrapAroundIndex(i+BoardWidth)] +
			b.arr[wrapAroundIndex(i+BoardWidth+1)]

		if sum < 2 || sum > 3 {
			tmp[i] = CellDead
		} else if b.arr[i] == CellDead && sum == 3 {
			tmp[i] = CellLive
		}
	}
	b.arr = tmp
	return nil
}
