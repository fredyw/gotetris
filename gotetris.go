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
// OUT OF OR I

package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
)

const (
	author string = "Fredy Wijaya"
	leftX  int    = 1
	leftY  int    = 1
	rightX int    = 21
	rightY int    = 20
	xStep  int    = 2
	yStep  int    = 1
)

type coordinate struct {
	y      int
	x      int
	filled bool
}

type game struct {
	coordinates [][]coordinate
}

func (g *game) moveLeft() {
	for row := 0; row < len(g.coordinates); row++ {
		for col := 0; col < len(g.coordinates[row]); col++ {
			x := g.coordinates[row][col].x
			if x-xStep > leftX {
				g.coordinates[row][col].x -= xStep
			} else {
				break
			}
		}
	}
}

func (g *game) moveRight() {
	for row := 0; row < len(g.coordinates); row++ {
		for col := 0; col < len(g.coordinates[row]); col++ {
			x := g.coordinates[row][col].x
			if x+xStep < rightX-1 {
				g.coordinates[row][col].x += xStep
			} else {
				break
			}
		}
	}
}

func (g *game) moveDown() {
	for row := 0; row < len(g.coordinates); row++ {
		for col := 0; col < len(g.coordinates[row]); col++ {
			y := g.coordinates[row][col].y
			if y+yStep < rightY {
				g.coordinates[row][col].y += yStep
			} else {
				break
			}
		}
	}
}

func (g *game) rotate() {
	// transpose
	tmpCoordinates := [][]coordinate{}
	for row := 0; row < len(g.coordinates); row++ {
		tmpCoordinates = append(tmpCoordinates, []coordinate{})
		for col := 0; col < len(g.coordinates[row]); col++ {
			tmpCoordinates[row] = append(tmpCoordinates[row], g.coordinates[col][row])
		}
	}

	for row := 0; row < len(g.coordinates); row++ {
		for col := 0; col < len(g.coordinates[row]); col++ {
			g.coordinates[row][col].filled = tmpCoordinates[row][col].filled
		}
	}

	// reverse
	for row := 0; row < len(g.coordinates); row++ {
		lcol := 0
		rcol := len(g.coordinates[row]) - 1
		for lcol < len(g.coordinates[row])/2 {
			tmp := g.coordinates[row][rcol].filled
			g.coordinates[row][rcol].filled = g.coordinates[row][lcol].filled
			g.coordinates[row][lcol].filled = tmp
			lcol++
			rcol--
		}
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

func drawRightLine() {
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
		termbox.SetCell(i, rightY+1, c, colorDefault, colorDefault)
	}
}

func drawBottomLine() {
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
	for row := 0; row < len(g.coordinates); row++ {
		step := 0
		for col := 0; col < len(g.coordinates[row]); col++ {
			c := '*'
			filled := g.coordinates[row][col].filled
			if !filled {
				c = ' '
			}
			x := g.coordinates[row][col].x + step
			y := g.coordinates[row][col].y
			step++
			termbox.SetCell(x, y, c, colorDefault, colorDefault)
		}
	}
}

func redrawAll(game *game) {
	colorDefault := termbox.ColorDefault
	termbox.Clear(colorDefault, colorDefault)

	drawBox()
	drawBlock(game)

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

	game := &game{
		coordinates: [][]coordinate{
			{
				{4, 4, false}, {4, 5, false}, {4, 6, false},
			},
			{
				{5, 4, true}, {5, 5, true}, {5, 6, true},
			},
			{
				{6, 4, false}, {6, 5, false}, {6, 6, true},
			},
		},
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
