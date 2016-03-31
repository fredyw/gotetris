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

type grid struct {
	leftX  int
	leftY  int
	rightX int
	rightY int
}

const (
	author    string = "Fredy Wijaya"
	xStep     int    = 2
	yStep     int    = 1
	numShapes int32  = 7
	maxRow    int    = 25
	maxCol    int    = 25
)

type block [][]coordinate

var (
	leftGrid grid = grid{
		leftX:  1,
		leftY:  0,
		rightX: 22,
		rightY: 20,
	}
	rightGrid grid = grid{
		leftX:  leftGrid.rightX,
		leftY:  leftGrid.leftY,
		rightX: 46,
		rightY: leftGrid.rightY,
	}
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
	newBlock block
	block    block
}

func (g *game) moveLeft() {
	collision := false
	for row := 0; row < len(g.newBlock); row++ {
		for col := 0; col < len(g.newBlock[row]); col++ {
			g.newBlock[row][col].x -= xStep
			x := g.newBlock[row][col].x
			y := g.newBlock[row][col].y
			if x >= 0 && y >= 0 {
				if (g.newBlock[row][col].x <= leftGrid.leftX && g.newBlock[row][col].filled) ||
					(g.block[y][x].filled && g.block[y][x].filled == g.newBlock[row][col].filled) {
					collision = true
				}
			}
		}
	}
	if collision {
		for row := 0; row < len(g.newBlock); row++ {
			for col := 0; col < len(g.newBlock[row]); col++ {
				g.newBlock[row][col].x += xStep
			}
		}
	}
}

func (g *game) moveRight() {
	collision := false
	for row := 0; row < len(g.newBlock); row++ {
		for col := len(g.newBlock[row]) - 1; col >= 0; col-- {
			g.newBlock[row][col].x += xStep
			x := g.newBlock[row][col].x
			y := g.newBlock[row][col].y
			if (g.newBlock[row][col].x+1 >= leftGrid.rightX && g.newBlock[row][col].filled) ||
				(g.block[y][x].filled && g.block[y][x].filled == g.newBlock[row][col].filled) {
				collision = true
			}
		}
	}
	if collision {
		for row := 0; row < len(g.newBlock); row++ {
			for col := 0; col < len(g.newBlock[row]); col++ {
				g.newBlock[row][col].x -= xStep
			}
		}
	}
}

func (g *game) moveDown() {
	// check for collision
	collision := false
	for row := 0; row < len(g.newBlock); row++ {
		for col := 0; col < len(g.newBlock[row]); col++ {
			g.newBlock[row][col].y += yStep
			if g.newBlock[row][col].y >= leftGrid.rightY && g.newBlock[row][col].filled {
				collision = true
			} else {
				x := g.newBlock[row][col].x
				y := g.newBlock[row][col].y
				if x >= 0 && g.block[y][x].filled && g.block[y][x].filled == g.newBlock[row][col].filled {
					collision = true
				}
			}
		}
	}
	if collision {
		for row := 0; row < len(g.newBlock); row++ {
			for col := 0; col < len(g.newBlock[row]); col++ {
				g.newBlock[row][col].y -= yStep
			}
		}

		for row := 0; row < len(g.newBlock); row++ {
			for col := 0; col < len(g.newBlock[row]); col++ {
				x := g.newBlock[row][col].x
				y := g.newBlock[row][col].y
				filled := g.newBlock[row][col].filled
				if filled {
					g.block[y][x].filled = filled
					g.block[y][x+1].filled = filled
				}
			}
		}
		removeBlock(g)
		g.newBlock = createNewBlock()
	}
}

func (g *game) rotate() {
	// keep a backup for reverting
	oldBlock := block{}
	for row := 0; row < len(g.newBlock); row++ {
		oldBlock = append(oldBlock, []coordinate{})
		for col := 0; col < len(g.newBlock[row]); col++ {
			oldCoordinate := coordinate{
				x:      g.newBlock[row][col].x,
				y:      g.newBlock[row][col].y,
				filled: g.newBlock[row][col].filled,
			}
			oldBlock[row] = append(oldBlock[row], oldCoordinate)
		}
	}

	// transpose
	tmpBlock := block{}
	for row := 0; row < len(g.newBlock); row++ {
		tmpBlock = append(tmpBlock, []coordinate{})
		for col := 0; col < len(g.newBlock[row]); col++ {
			tmpBlock[row] = append(tmpBlock[row], g.newBlock[col][row])
		}
	}

	for row := 0; row < len(g.newBlock); row++ {
		for col := 0; col < len(g.newBlock[row]); col++ {
			g.newBlock[row][col].filled = tmpBlock[row][col].filled
		}
	}

	// reverse
	for row := 0; row < len(g.newBlock); row++ {
		lcol := 0
		rcol := len(g.newBlock[row]) - 1
		for lcol < len(g.newBlock[row])/2 {
			tmp := g.newBlock[row][rcol].filled
			g.newBlock[row][rcol].filled = g.newBlock[row][lcol].filled
			g.newBlock[row][lcol].filled = tmp
			lcol++
			rcol--
		}
	}

	collision := false
	for row := 0; row < len(g.newBlock); row++ {
		for col := len(g.newBlock[row]) - 1; col >= 0; col-- {
			x := g.newBlock[row][col].x
			y := g.newBlock[row][col].y
			if g.newBlock[row][col].x+1 >= leftGrid.rightX && g.newBlock[row][col].filled ||
				g.newBlock[row][col].x <= leftGrid.leftX && g.newBlock[row][col].filled ||
				g.newBlock[row][col].y >= leftGrid.rightY && g.newBlock[row][col].filled ||
				(g.block[y][x].filled && g.block[y][x].filled == g.newBlock[row][col].filled) {
				collision = true
			}
		}
	}
	if collision {
		g.newBlock = oldBlock
	}
	removeBlock(g)
}

func removeBlock(g *game) {
	for true {
		rows := []int{}
		for row := leftGrid.leftY + 1; row < leftGrid.rightY; row++ {
			allFilled := true
			for col := leftGrid.leftX + 1; col < leftGrid.rightX; col++ {
				filled := g.block[row][col].filled
				if !filled {
					allFilled = false
				}
			}
			if allFilled {
				rows = append(rows, row)
			}
		}
		if len(rows) > 0 {
			lastRow := rows[len(rows)-1]
			for row := lastRow; row > 0; row-- {
				for col := 0; col < len(g.block[row]); col++ {
					g.block[row][col].filled = g.block[row-1][col].filled
				}
			}
		} else {
			break
		}
	}
}

func drawNewBlock(g *game) {
	colorDefault := termbox.ColorDefault
	for row := 0; row < len(g.newBlock); row++ {
		for col := 0; col < len(g.newBlock[row]); col++ {
			c := '\u2588'
			x := g.newBlock[row][col].x
			y := g.newBlock[row][col].y
			filled := g.newBlock[row][col].filled
			if !filled {
				if x >= 0 {
					if !g.block[y][x].filled {
						c = ' '
					}
				}
			}
			termbox.SetCell(x, y, c, colorDefault, colorDefault)
			termbox.SetCell(x+1, y, c, colorDefault, colorDefault)
		}
	}
}

func drawBlock(g *game) {
	colorDefault := termbox.ColorDefault
	for row := 0; row < len(g.block); row++ {
		for col := 0; col < len(g.block[row]); col++ {
			c := '\u2588'
			x := g.block[row][col].x
			y := g.block[row][col].y
			filled := g.block[row][col].filled
			if !filled {
				c = ' '
			}
			termbox.SetCell(x, y, c, colorDefault, colorDefault)
		}
	}
}

func drawText(x, y int, text string) {
	colorDefault := termbox.ColorDefault
	for _, ch := range text {
		termbox.SetCell(x, y, ch, colorDefault, colorDefault)
		x++
	}
}

func drawLeftGridTopLine() {
	colorDefault := termbox.ColorDefault
	for i := leftGrid.leftX; i <= leftGrid.rightX; i++ {
		var c rune
		if i == leftGrid.leftX {
			c = '\u250c'
		} else if i == leftGrid.rightX {
			c = '\u2510'
		} else {
			c = '\u2500'
		}
		termbox.SetCell(i, leftGrid.leftY, c, colorDefault, colorDefault)
	}
}

func drawLeftGridLeftLine() {
	colorDefault := termbox.ColorDefault
	for i := leftGrid.leftY + 1; i <= leftGrid.rightY; i++ {
		c := '\u2502'
		termbox.SetCell(leftGrid.leftX, i, c, colorDefault, colorDefault)
	}
}

func drawLeftGridRightLine() {
	colorDefault := termbox.ColorDefault
	for i := leftGrid.leftY + 1; i <= leftGrid.rightY; i++ {
		c := '\u2502'
		termbox.SetCell(leftGrid.rightX, i, c, colorDefault, colorDefault)
	}
}

func drawLeftGridBottomLine() {
	colorDefault := termbox.ColorDefault
	for i := leftGrid.leftX; i <= leftGrid.rightX; i++ {
		var c rune
		if i == leftGrid.leftX {
			c = '\u2514'
		} else if i == leftGrid.rightX {
			c = '\u2518'
		} else {
			c = '\u2500'
		}
		termbox.SetCell(i, leftGrid.rightY, c, colorDefault, colorDefault)
	}
}

func drawLeftGrid() {
	drawLeftGridTopLine()
	drawLeftGridLeftLine()
	drawLeftGridRightLine()
	drawLeftGridBottomLine()
}

func drawRightGridTopLine() {
	colorDefault := termbox.ColorDefault
	for i := rightGrid.leftX; i <= rightGrid.rightX; i++ {
		var c rune
		if i == rightGrid.leftX {
			c = '\u252c'
		} else if i == rightGrid.rightX {
			c = '\u2510'
		} else {
			c = '\u2500'
		}
		termbox.SetCell(i, rightGrid.leftY, c, colorDefault, colorDefault)
	}
}

func drawRightGridRightLine() {
	colorDefault := termbox.ColorDefault
	for i := rightGrid.leftY + 1; i <= rightGrid.rightY; i++ {
		c := '\u2502'
		termbox.SetCell(rightGrid.rightX, i, c, colorDefault, colorDefault)
	}
}

func drawRightGridBottomLine() {
	colorDefault := termbox.ColorDefault
	for i := rightGrid.leftX; i <= rightGrid.rightX; i++ {
		var c rune
		if i == rightGrid.leftX {
			c = '\u2534'
		} else if i == rightGrid.rightX {
			c = '\u2518'
		} else {
			c = '\u2500'
		}
		termbox.SetCell(i, rightGrid.rightY, c, colorDefault, colorDefault)
	}
}

func drawLevel(level int) {
	x := rightGrid.leftX + 2
	drawText(x, 1, fmt.Sprintf("Level : %d", level))
}

func drawScore(score int) {
	x := rightGrid.leftX + 2
	drawText(x, 2, fmt.Sprintf("Score : %d", score))
}

func drawSeparator1() {
	colorDefault := termbox.ColorDefault
	for i := rightGrid.leftX; i <= rightGrid.rightX; i++ {
		var c rune
		if i == rightGrid.leftX {
			c = '\u251C'
		} else if i == rightGrid.rightX {
			c = '\u2524'
		} else {
			c = '\u2500'
		}
		termbox.SetCell(i, 3, c, colorDefault, colorDefault)
	}
}

func drawControls() {
	x := rightGrid.leftX + 2
	drawText(x, 4, "Move left  : \u2190")
	drawText(x, 5, "Move right : \u2192")
	drawText(x, 6, "Move down  : \u2193")
	drawText(x, 7, "Rotate     : Spacebar")
	drawText(x, 8, "Exit       : Esc")
}

func drawSeparator2() {
	colorDefault := termbox.ColorDefault
	for i := rightGrid.leftX; i <= rightGrid.rightX; i++ {
		var c rune
		if i == rightGrid.leftX {
			c = '\u251C'
		} else if i == rightGrid.rightX {
			c = '\u2524'
		} else {
			c = '\u2500'
		}
		termbox.SetCell(i, 18, c, colorDefault, colorDefault)
	}
}

func drawAuthor() {
	x := rightGrid.leftX + 2
	drawText(x, 19, fmt.Sprintf("Author : %s", author))
}

func drawRightGrid() {
	drawRightGridTopLine()
	drawRightGridRightLine()
	drawRightGridBottomLine()

	drawLevel(1)
	drawScore(0)
	drawSeparator1()
	drawControls()
	drawSeparator2()
	drawAuthor()
}

func redrawAll(game *game) {
	colorDefault := termbox.ColorDefault
	termbox.Clear(colorDefault, colorDefault)

	drawBlock(game)
	drawNewBlock(game)

	drawLeftGrid()
	drawRightGrid()

	termbox.Flush()
}

func createNewBlock() block {
	shape := shapes[rand.Int31n(numShapes)]
	// create a copy
	newBlock := block{}
	for row := 0; row < len(shape); row++ {
		newBlock = append(newBlock, []coordinate{})
		for col := 0; col < len(shape[row]); col++ {
			newBlock[row] = append(newBlock[row], coordinate{})
			newBlock[row][col].x = shape[row][col].x
			newBlock[row][col].y = shape[row][col].y
			newBlock[row][col].filled = shape[row][col].filled
		}
	}
	return newBlock
}

func initBlock() block {
	block := block{}
	for row := 0; row <= maxRow; row++ {
		block = append(block, []coordinate{})
		for col := 0; col <= maxCol; col++ {
			block[row] = append(block[row], coordinate{
				x:      col,
				y:      row,
				filled: false,
			})
		}
	}
	return block
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
		newBlock: createNewBlock(),
		block:    initBlock(),
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
