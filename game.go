package battleships

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

const (
	// Rows defines number of rows of the game's board
	Rows = 10
	// Cols defines number of cols of the game's board
	Cols = 10

	inputRegex          = "^[A-J](10|[1-9])$"
	horizontalDirection = 0
	verticalDirection   = 1

	// EmptySlot defines a field, that doesn't contain any ship and was hit hit so far
	EmptySlot = '-'
	// ShipSlot defines a field, that contains am undamaged ship
	ShipSlot = 'S'
	// HitShipSlot defines a field with a ship, that is already hit
	HitShipSlot = 'X'
	// MissedSlot defines a field without a ship, but that was shot at already
	MissedSlot = 'O'
)

// PatternMismatch defines error used, when there is not match with the required pattern
type PatternMismatch struct {
	input string
}

func (e PatternMismatch) Error() string {
	return fmt.Sprintf("%v doesn't match the pattern %v", e.input, inputRegex)
}

// Board describes a game board used to store information about the current state of a game.
type Board [Rows][Cols]byte

// At is a convenient method used to access board field using the Position object. Returns byte at a specified location in the board
func (b *Board) At(p Position) byte {
	return b[p.row][p.col]
}

// Set is a convenient method used to set a new value in the board indexing it with a Position object
func (b *Board) Set(p Position, val byte) {
	b[p.row][p.col] = val
}

// Ship describes a single ship object used in the game
type Ship struct {
	size   uint8
	health uint8
}

func (s *Ship) hit() bool {
	s.health--
	return s.health == 0
}

// NewShip creates a new ship with given size and full health
func NewShip(size uint8) Ship {
	return Ship{
		size:   size,
		health: size,
	}
}

// Game defines an object used to initialize and start a new game
type Game struct {
	Stats Statistics

	shipsData   map[Position]*Ship
	board       Board
	initialized bool
}

// Statistics defines information about current state of the game
type Statistics struct {
	ShotsFired   int
	InitialShips int
	SunkShips    int
}

// Position describes indexes used to access game's board
type Position struct {
	row, col uint8
}

// Shot method allows to try to hit a ship at given position.
// First returned value is true, if a ship was hit. At the same time, if it was the last slot of a ship, true will be returned as second value
// Method returns error, if called before the game is iniatialized
func (g *Game) Shot(pos Position) (bool, bool, error) {
	if !g.initialized {
		return false, false, errors.New("Game not initialized")
	}
	g.Stats.ShotsFired++

	if g.board.At(pos) == ShipSlot {
		g.board.Set(pos, HitShipSlot)
		s := g.shipsData[pos]
		sunk := s.hit()
		if sunk {
			g.Stats.SunkShips++
		}
		return true, sunk, nil
	} else if g.board.At(pos) == EmptySlot {
		g.board.Set(pos, MissedSlot)
	}
	return false, false, nil
}

// FillBoard fills randomly the game's board with given ships.
// After that, the game is fully initialized and ready to be played
func (g *Game) FillBoard(ships []Ship) {
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			g.board[i][j] = EmptySlot
		}
	}
	g.shipsData = make(map[Position]*Ship)
	g.Stats.InitialShips = len(ships)

	rand := rand.New(rand.NewSource(time.Now().Unix()))

	for _, s := range ships {
		placed := false
		tries := 0
		for !placed {
			tries++
			direction := rand.Intn(2)
			maxRow := Rows
			maxCol := Cols
			if direction == horizontalDirection {
				maxRow = Rows - int(s.size) + 1
			} else {
				maxCol = Cols - int(s.size) + 1
			}
			pos := randomPosition(rand, maxRow, maxCol)

			if canPlaceShip(g, s, pos, direction) {
				placeShip(g, s, pos, direction)
				placed = true
			}
			if tries == 50 {
				return
			}
		}
	}
	g.initialized = true
}

// Playable returns true, if there are still ships alive in the current game
func (g *Game) Playable() bool {
	return g.initialized && g.Stats.SunkShips < g.Stats.InitialShips
}

// Board returns deep copy of a game's board. Parametr describes, if ships will be marked on the board or not
func (g *Game) Board(hiddenShips bool) *Board {
	b := &Board{}
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			if hiddenShips && g.board[i][j] == ShipSlot {
				b[i][j] = EmptySlot
			} else {
				b[i][j] = g.board[i][j]
			}
		}
	}

	return b
}

func randomPosition(rand *rand.Rand, maxR, maxC int) Position {
	row := rand.Intn(maxR)
	col := rand.Intn(maxC)

	return Position{row: uint8(row), col: uint8(col)}
}

func canPlaceShip(g *Game, ship Ship, pos Position, direction int) bool {

	for i := uint8(0); i < ship.size; i++ {
		switch direction {
		case horizontalDirection:
			if !isValidPosition(g, pos.row, pos.col+i) {
				return false
			}
		case verticalDirection:
			if !isValidPosition(g, pos.row+i, pos.col) {
				return false
			}
		}
	}
	return true
}

func isValidPosition(g *Game, row, col uint8) bool {
	return isWithinBoard(row, col) && !isAnotherShipInNeighbourhood(g, row, col)
}

func isWithinBoard(row, col uint8) bool {
	return row >= 0 && row < Rows && col >= 0 && col < Cols
}

func isAnotherShipInNeighbourhood(g *Game, row, col uint8) bool {
	minR := max(0, int8(row-1))
	maxR := min(Rows-1, int8(row+1))
	minC := max(0, int8(col-1))
	maxC := min(Cols-1, int8(col+1))

	for i := minR; i <= maxR; i++ {
		for j := minC; j <= maxC; j++ {
			if g.board[i][j] == ShipSlot {
				return true
			}
		}
	}
	return false
}

func min(x, y int8) int8 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int8) int8 {
	if x > y {
		return x
	}
	return y
}

func placeShip(g *Game, ship Ship, pos Position, direction int) {
	for i := uint8(0); i < ship.size; i++ {
		switch direction {
		case horizontalDirection:
			g.addShip(&ship, Position{row: pos.row, col: pos.col + i})
		case verticalDirection:
			g.addShip(&ship, Position{row: pos.row + i, col: pos.col})
		}
	}
}

func (g *Game) addShip(ship *Ship, pos Position) {
	g.board.Set(pos, ShipSlot)
	g.shipsData[pos] = ship
}

// ConvertInputToPosition allows to convert text input in form [A-Z][1-10] to corresponding (row,column) position.
// Returns error if the input doesn't match required pattern
func ConvertInputToPosition(input string) (*Position, error) {
	matched, err := regexp.MatchString(inputRegex, input)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, PatternMismatch{input}
	}

	letter := input[0]
	number := input[1:]
	row := letter - 'A'
	col, _ := strconv.ParseUint(number, 10, 8)

	return &Position{row: row, col: uint8(col - 1)}, nil
}
