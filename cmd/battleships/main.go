package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"bytes"

	"github.com/jkosecki/battleships"
)

const (
	inputMessage     = "Please insert a new position in form '[A-J][1-10]': "
	readErrorMessage = "Problems while reading input occured. Please try again"
	errorMessage     = "Your input '%v' doesn't match the required form. Please type again: "
)

func main() {

	p := newPlayer(os.Stdin, os.Stdout)
	g := &battleships.Game{}

	g.FillBoard([]battleships.Ship{
		battleships.NewShip(5),
		battleships.NewShip(4),
		battleships.NewShip(4),
	})

	for g.Playable() {
		printBoard(g.Board(true))
		pos := p.GetShotPosition()
		hit, sunk, err := g.Shot(pos)
		if err != nil {
			fmt.Println(err)
			return
		}
		if hit {
			fmt.Println("\nYou've hit a ship")
		}
		if sunk {
			ships := g.Stats.InitialShips
			fmt.Printf("A ship has sunk! %v/%v still alive\n", ships-g.Stats.SunkShips, ships)
		}
		fmt.Println()
	}
	printBoard(g.Board(false))
	fmt.Printf("Game over. All ships are sunk after %v shots\n", g.Stats.ShotsFired)
}

type consolePlayer struct {
	in  bufio.Scanner
	out io.Writer
}

func newPlayer(r io.Reader, w io.Writer) consolePlayer {
	return consolePlayer{
		in:  *bufio.NewScanner(r),
		out: w,
	}
}

func (p *consolePlayer) GetShotPosition() battleships.Position {
	fmt.Fprint(p.out, inputMessage)

	for {
		p.in.Scan()
		input := strings.ToUpper(p.in.Text())

		pos, err := battleships.ConvertInputToPosition(input)

		if err != nil {
			fmt.Fprintf(p.out, errorMessage, input)
		} else {
			return *pos
		}
	}
}

func printBoard(board *battleships.Board) {
	buf := bytes.Buffer{}
	buf.WriteString("  ")
	for i := 0; i < battleships.Cols; i++ {
		buf.WriteString(fmt.Sprintf("%3d", i+1))
	}
	buf.WriteString("\n")

	for i := 0; i < battleships.Rows; i++ {
		buf.WriteString(fmt.Sprintf("%2c", 'A'+i))
		for j := 0; j < battleships.Cols; j++ {
			buf.WriteString(fmt.Sprintf("%3c", board[i][j]))
		}
		buf.WriteString("\n")
	}
	fmt.Print(buf.String())
}
