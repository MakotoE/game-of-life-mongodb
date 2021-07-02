package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

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
	board, err := NewDatabaseBoard()
	require.Nil(t, err)
	defer board.Close()

	for i, test := range tests {
		initial := make(map[[2]int]Cell, BoardWidth*BoardHeight)
		for x := 0; x < BoardWidth; x++ {
			for y := 0; y < BoardHeight; y++ {
				initial[[2]int{x, y}] = CellDead
			}
		}

		for _, coordinate := range test.initiallyActivatedCells {
			initial[coordinate] = CellLive
		}

		for coordinate, cell := range initial {
			require.Nil(t, board.Set(coordinate[0], coordinate[1], cell), i)
		}

		require.Nil(t, board.Tick(), i)

		expected := make(map[[2]int]Cell, BoardWidth*BoardHeight)
		for x := 0; x < BoardWidth; x++ {
			for y := 0; y < BoardHeight; y++ {
				expected[[2]int{x, y}] = CellDead
			}
		}

		for _, coordinate := range test.expectedLive {
			expected[coordinate] = CellLive
		}

		for coordinate, expectedCell := range expected {
			result, err := board.Cell(coordinate[0], coordinate[1])
			require.Nil(t, err, i)
			assert.Equal(t, expectedCell, result, i)
		}
	}
}

func TestGetCellWrapAround(t *testing.T) {
	arr := [BoardWidth * BoardHeight]Cell{}
	arr[BoardWidth*BoardHeight-1] = CellLive
	assert.Equal(t, CellDead, arr[wrapAroundIndex(0)])
	assert.Equal(t, CellLive, arr[wrapAroundIndex(-1)])          // Left
	assert.Equal(t, CellDead, arr[wrapAroundIndex(-BoardWidth)]) // Up

	arr[BoardWidth*BoardHeight-BoardWidth-1] = CellLive
	assert.Equal(t, CellDead, arr[wrapAroundIndex(-BoardWidth)])
}
