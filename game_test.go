package battleships

import (
	"testing"
)

func TestConvertInputToPosition_noMatch(t *testing.T) {
	inputs := []string{"A0", "A11", "K1"}

	for _, input := range inputs {
		_, err := ConvertInputToPosition(input)

		_, ok := err.(PatternMismatch)

		if !ok {
			t.Errorf("No PatternMismatch for input %v - err: %v", input, err)
		}
	}
}

func TestConvertInputToPosition_positionsReturned(t *testing.T) {
	data := []struct {
		in  string
		out Position
	}{
		{"A1", Position{0, 0}},
		{"A10", Position{0, 9}},
		{"J1", Position{9, 0}},
		{"J10", Position{9, 9}},
		{"B5", Position{1, 4}},
		{"C6", Position{2, 5}},
		{"D7", Position{3, 6}},
		{"E8", Position{4, 7}},
	}

	for _, d := range data {
		pos, err := ConvertInputToPosition(d.in)

		if err != nil {
			t.Errorf("Error has been thrown %v", err)
		}

		if err == nil && *pos != d.out {
			t.Errorf("Expected: %v, received: %v", d.out, *pos)
		}
	}
}

func TestNewShip_helthAndSizeTheSame(t *testing.T) {
	data := []uint8{1, 2, 3, 4, 5}

	for _, d := range data {
		s := NewShip(d)

		if s.health != s.size {
			t.Errorf("Health (%v) and size (%v) do not match", s.health, s.size)
		}
	}
}

func TestAt_correctValueReturned(t *testing.T) {
	data := []struct {
		row, col uint8
		val      byte
	}{
		{0, 0, 'X'},
		{5, 5, '-'},
		{4, 2, 'S'},
		{9, 0, 'o'},
	}

	b := Board{}

	for _, d := range data {
		b[d.row][d.col] = d.val

		got := b.At(Position{d.row, d.col})
		if d.val != got {
			t.Errorf("Expected value: %v, got: %v", d.val, got)
		}
	}
}

func TestSet_correctValueSet(t *testing.T) {
	data := []struct {
		row, col uint8
		val      byte
	}{
		{0, 0, 'X'},
		{5, 5, '-'},
		{4, 2, 'S'},
		{9, 0, 'o'},
	}

	b := Board{}

	for _, d := range data {
		b.Set(Position{d.row, d.col}, d.val)

		got := b[d.row][d.col]
		if d.val != got {
			t.Errorf("Expected value: %v, got: %v", d.val, got)
		}
	}
}

func TestFillBoard_boardFilledWithShips(t *testing.T) {
	g := Game{}
	ships := []Ship{NewShip(5), NewShip(4)}

	g.FillBoard(ships)

	if !g.initialized {
		t.Error("Game has been not initialized")
	}
	if g.Stats.InitialShips != len(ships) {
		t.Error("Statistic of initial ships not filled")
	}

	var shipSlots uint8
	var emptySlots uint8
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			if g.board[i][j] == ShipSlot {
				shipSlots++
			} else {
				emptySlots++
			}
		}
	}
	var expectedShipSlots uint8
	for _, s := range ships {
		expectedShipSlots += s.size
	}
	expectedEmptySlots := Rows*Cols - expectedShipSlots

	if expectedShipSlots != shipSlots {
		t.Errorf("Expected number of ship slots: %v, got: %v", expectedShipSlots, shipSlots)
	}
	if expectedEmptySlots != emptySlots {
		t.Errorf("Expected number of empty slots: %v, got: %v", expectedEmptySlots, emptySlots)
	}
}

func TestPlayable(t *testing.T) {
	data := []struct {
		initialized         bool
		iniatialShips, sunk int
		expected            bool
	}{
		{false, 0, 0, false},
		{false, 1, 0, false},
		{true, 1, 1, false},
		{true, 1, 0, true},
	}

	g := Game{}
	for _, d := range data {
		g.initialized = d.initialized
		g.Stats.InitialShips = d.iniatialShips
		g.Stats.SunkShips = d.sunk

		if g.Playable() != d.expected {
			t.Errorf("Expected status: %v, got: %v for data: %v", d.expected, g.Playable(), d)
		}
	}
}
