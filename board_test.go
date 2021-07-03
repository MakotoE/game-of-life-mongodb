package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
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

	arrayBoard := NewArrayBoard()
	databaseBoard, err := NewDatabaseBoard()
	require.Nil(t, err)
	defer databaseBoard.Close()

	for _, board := range [...]Board{&arrayBoard, &databaseBoard} {
		for i, test := range tests {
			initialCells := NewCellArray(test.initiallyActivatedCells)

			for i, cell := range initialCells.Arr {
				err := board.Set(i%BoardWidth, i/BoardWidth, cell)
				require.Nil(t, err, "%v %v", reflect.TypeOf(board), i)
			}

			require.Nil(t, board.Tick(), "%v %v", reflect.TypeOf(board), i)

			result, err := board.Cells()
			require.Nil(t, err, "%v %v", reflect.TypeOf(board), i)

			expected := NewCellArray(test.expectedLive)
			assert.Equal(t, expected, result, "%v %v", reflect.TypeOf(board), i)
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
