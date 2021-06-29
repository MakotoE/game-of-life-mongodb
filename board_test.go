package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func activateCell(board *[BoardWidth * BoardHeight]Cell, x int, y int) {
	board[y*BoardWidth+x] = CellLive
}

func TestBoard_Tick(t *testing.T) {
	tests := []struct {
		initiallyActivatedCells [][2]int
		expectedLive            [][2]int
	}{
		// Rule 1
		{
			nil,
			nil,
		},
		{
			[][2]int{{0, 0}},
			nil,
		},
		{
			[][2]int{{0, 0}, {1, 0}},
			nil,
		},
		// Rule 2
		{
			[][2]int{
				{0, 0},
				{1, 0},
				{0, 1},
			},
			[][2]int{
				{0, 0},
				{1, 0},
				{0, 1},
				{1, 1},
			},
		},
		{
			[][2]int{
				{0, 0},
				{1, 0},
				{2, 0},
			},
			[][2]int{
				{1, 0},
				{1, 1},
				{1, BoardHeight - 1},
			},
		},
		{
			[][2]int{
				{0, 0},
				{1, 0},
				{0, 1},
				{1, 1},
			},
			[][2]int{
				{0, 0},
				{1, 0},
				{0, 1},
				{1, 1},
			},
		},
		// Rule 3 and 4
		{
			[][2]int{
				{1, 0},
				{0, 1},
				{1, 1},
				{2, 1},
				{1, 2},
			},
			[][2]int{
				{0, 0},
				{1, 0},
				{2, 0},
				{0, 1},
				{2, 1},
				{0, 2},
				{1, 2},
				{2, 2},
			},
		},
	}

	for i, test := range tests {
		board := Board{}
		for _, coordinate := range test.initiallyActivatedCells {
			activateCell(&board.arr, coordinate[0], coordinate[1])
		}

		board.Tick()

		expected := [BoardWidth * BoardHeight]Cell{}
		for _, coordinate := range test.expectedLive {
			activateCell(&expected, coordinate[0], coordinate[1])
		}

		assert.Equal(t, expected, board.arr, i)
	}
}

func TestGetCellWrapAround(t *testing.T) {
	arr := [BoardWidth * BoardHeight]Cell{}
	arr[BoardWidth*BoardHeight-1] = CellLive
	assert.Equal(t, int(CellDead), getCellWrapAround(&arr, 0))
	assert.Equal(t, int(CellLive), getCellWrapAround(&arr, -1))          // Left
	assert.Equal(t, int(CellDead), getCellWrapAround(&arr, -BoardWidth)) // Up

	arr[BoardWidth*BoardHeight-BoardWidth-1] = CellLive
	assert.Equal(t, int(CellDead), getCellWrapAround(&arr, -BoardWidth))
}
