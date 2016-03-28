// The MIT License (MIT)
//
// Copyright (c) 2016 Fredy Wijaya
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"time"
)

const (
	author string = "Fredy Wijaya"
	leftX  int    = 2
	leftY  int    = 0
	rightX int    = 22
	rightY int    = 20
	xStep  int    = 1
	yStep  int    = 1
)

type block [][]coordinate

var (
	shapes []block = []block{
		{
			{
				{1, 8, false}, {1, 10, false}, {1, 12, false}, {1, 14, false},
			},
			{
				{2, 8, true}, {2, 10, true}, {2, 12, true}, {2, 14, true},
			},
			{
				{3, 8, false}, {3, 10, false}, {3, 12, false}, {3, 14, false},
			},
			{
				{4, 8, false}, {4, 10, false}, {4, 12, false}, {4, 14, false},
			},
		},
		{
			{
				{1, 8, true}, {1, 10, false}, {1, 12, false},
			},
			{
				{2, 8, true}, {2, 10, true}, {2, 12, true},
			},
			{
				{3, 8, false}, {3, 10, false}, {3, 12, false},
			},
		},
		{
			{
				{1, 8, false}, {1, 10, false}, {1, 12, true},
			},
			{
				{2, 8, true}, {2, 10, true}, {2, 12, true},
			},
			{
				{3, 8, false}, {3, 10, false}, {3, 12, false},
			},
		},
		{
			{
				{1, 10, true}, {1, 12, true},
			},
			{
				{2, 10, true}, {2, 12, true},
			},
		},
		{
			{
				{1, 8, false}, {1, 10, true}, {1, 12, true},
			},
			{
				{2, 8, true}, {2, 10, true}, {2, 12, false},
			},
			{
				{3, 8, false}, {3, 10, false}, {3, 12, false},
			},
		},
		{
			{
				{1, 8, false}, {1, 10, true}, {1, 12, false},
			},
			{
				{2, 8, true}, {2, 10, true}, {2, 12, true},
			},
			{
				{3, 8, false}, {3, 10, false}, {3, 12, false},
			},
		},
		{
			{
				{1, 8, true}, {1, 10, true}, {1, 12, false},
			},
			{
				{2, 8, false}, {2, 10, true}, {2, 12, true},
			},
			{
				{3, 8, false}, {3, 10, false}, {3, 12, false},
			},
		},
	}
)

type coordinate struct {
	y      int
	x      int
	filled bool
}

type game struct {
	block block
}

func (g *game) moveLeft() {
	revert := false
	for row := 0; row < len(g.block); row++ {
		for col := 0; col < len(g.block[row]); col++ {
			g.block[row][col].x -= xStep
			if g.block[row][col].x <= leftX && g.block[row][col].filled {
				revert = true
			}
		}
	}
	if revert {
		for row := 0; row < len(g.block); row++ {
			for col := 0; col < len(g.block[row]); col++ {
				g.block[row][col].x += xStep
			}
		}
	}
}

func (g *game) moveRight() {
	revert := false
	for row := 0; row < len(g.block); row++ {
		for col := len(g.block[row]) - 1; col >= 0; col-- {
			g.block[row][col].x += xStep
			if g.block[row][col].x+1 >= rightX && g.block[row][col].filled {
				revert = true
			}
		}
	}
	if revert {
		for row := 0; row < len(g.block); row++ {
			for col := 0; col < len(g.block[row]); col++ {
				g.block[row][col].x -= xStep
			}
		}
	}
}

func (g *game) moveDown() {
	revert := false
	for row := 0; row < len(g.block); row++ {
		for col := 0; col < len(g.block[row]); col++ {
			g.block[row][col].y += yStep
			if g.block[row][col].y >= rightY && g.block[row][col].filled {
				revert = true
			}
		}
	}
	if revert {
		for row := 0; row < len(g.block); row++ {
			for col := 0; col < len(g.block[row]); col++ {
				g.block[row][col].y -= yStep
			}
		}
	}
}

func (g *game) rotate() {
	// keep a backup for reverting
	oldBlock := block{}
	for row := 0; row < len(g.block); row++ {
		oldBlock = append(oldBlock, []coordinate{})
		for col := 0; col < len(g.block[row]); col++ {
			oldCoordinate := coordinate{
				x:      g.block[row][col].x,
				y:      g.block[row][col].y,
				filled: g.block[row][col].filled,
			}
			oldBlock[row] = append(oldBlock[row], oldCoordinate)
		}
	}

	// transpose
	tmpBlock := block{}
	for row := 0; row < len(g.block); row++ {
		tmpBlock = append(tmpBlock, []coordinate{})
		for col := 0; col < len(g.block[row]); col++ {
			tmpBlock[row] = append(tmpBlock[row], g.block[col][row])
		}
	}

	for row := 0; row < len(g.block); row++ {
		for col := 0; col < len(g.block[row]); col++ {
			g.block[row][col].filled = tmpBlock[row][col].filled
		}
	}

	// reverse
	for row := 0; row < len(g.block); row++ {
		lcol := 0
		rcol := len(g.block[row]) - 1
		for lcol < len(g.block[row])/2 {
			tmp := g.block[row][rcol].filled
			g.block[row][rcol].filled = g.block[row][lcol].filled
			g.block[row][lcol].filled = tmp
			lcol++
			rcol--
		}
	}

	revert := false
	for row := 0; row < len(g.block); row++ {
		for col := len(g.block[row]) - 1; col >= 0; col-- {
			if g.block[row][col].x+1 >= rightX && g.block[row][col].filled ||
				g.block[row][col].x <= leftX && g.block[row][col].filled ||
				g.block[row][col].y >= rightY && g.block[row][col].filled {
				revert = true
			}
		}
	}
	if revert {
		g.block = oldBlock
	}
}

func drawTopLine() {
	colorDefault := termbox.ColorDefault
	for i := leftX; i <= rightX; i++ {
		var c rune
		if i == leftX {
			c = '\u250c'
		} else if i == rightX {
			c = '\u2510'
		} else {
			c = '\u2500'
		}
		termbox.SetCell(i, leftY, c, colorDefault, colorDefault)
	}
}

func drawLeftLine() {
	colorDefault := termbox.ColorDefault
	for i := leftY + 1; i <= rightY; i++ {
		c := '\u2502'
		termbox.SetCell(leftX, i, c, colorDefault, colorDefault)
	}
}

func drawBottomLine() {
	colorDefault := termbox.ColorDefault
	for i := leftX; i <= rightX; i++ {
		var c rune
		if i == leftX {
			c = '\u2514'
		} else if i == rightX {
			c = '\u2518'
		} else {
			c = '\u2500'
		}
		termbox.SetCell(i, rightY, c, colorDefault, colorDefault)
	}
}

func drawRightLine() {
	colorDefault := termbox.ColorDefault
	for i := leftY + 1; i <= rightY; i++ {
		c := '\u2502'
		termbox.SetCell(rightX, i, c, colorDefault, colorDefault)
	}
}

func drawBox() {
	drawTopLine()
	drawLeftLine()
	drawRightLine()
	drawBottomLine()
}

func drawBlock(g *game) {
	colorDefault := termbox.ColorDefault
	for row := 0; row < len(g.block); row++ {
		for col := 0; col < len(g.block[row]); col++ {
			c := '\u2588'
			filled := g.block[row][col].filled
			if !filled {
				c = ' '
			}
			x := g.block[row][col].x
			y := g.block[row][col].y
			termbox.SetCell(x, y, c, colorDefault, colorDefault)
			//if col != len(g.coordinates[row])-1 {
			termbox.SetCell(x+1, y, c, colorDefault, colorDefault)
			//}
		}
	}
}

func redrawAll(game *game) {
	colorDefault := termbox.ColorDefault
	termbox.Clear(colorDefault, colorDefault)

	drawBlock(game)
	drawBox()

	termbox.Flush()
}

func runGame() {
	err := termbox.Init()
	if err != nil {
		errorAndExit(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	rand.Seed(time.Now().UTC().UnixNano())

	game := &game{
		block: shapes[rand.Int31n(7)],
	}

	redrawAll(game)
exitGame:
	for {
		for {
			select {
			case ev := <-eventQueue:
				switch ev.Key {
				case termbox.KeyEsc:
					break exitGame
				case termbox.KeyArrowLeft:
					game.moveLeft()
				case termbox.KeyArrowRight:
					game.moveRight()
				case termbox.KeyArrowDown:
					game.moveDown()
				case termbox.KeySpace:
					game.rotate()
				}
			}
			redrawAll(game)
		}
	}
}

func errorAndExit(message interface{}) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	runGame()
}
