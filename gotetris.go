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
	rightX int    = 20
	rightY int    = 20
	xStep  int    = 2
	yStep  int    = 1
)

type coordinate struct {
	x int
	y int
}

type block struct {
	coordinates []coordinate
}

type game struct {
	block block
}

//public static char[][] rotate(char[][] block) {
//int rowSize = block.length;
//int colSize = block[0].length;
//
//// tranpose the matrix
//char[][] result = new char[colSize][rowSize];
//for (int i = 0; i < rowSize; i++) {
//for (int j = 0; j < colSize; j++) {
//result[j][i] = block[i][j];
//}
//}
//// reverse the rows
//for (int i = 0; i < result.length; i++) {
//for (int j = 0, k = result[i].length - 1; j < result[i].length / 2; j++, k--) {
//char tmp = result[i][k];
//result[i][k] = result[i][j];
//result[i][j] = tmp;
//}
//}
//
//return result;
//}

func (b *block) moveLeft() {
	for i, _ := range b.coordinates {
		if b.coordinates[i].x-xStep > leftX {
			b.coordinates[i].x -= xStep
		}
	}
}

func (b *block) moveRight() {
	for i, _ := range b.coordinates {
		if b.coordinates[i].x+xStep < rightX {
			b.coordinates[i].x += xStep
		}
	}
}

func (b *block) moveDown() {
	for i, _ := range b.coordinates {
		if b.coordinates[i].y+yStep < rightY {
			b.coordinates[i].y += yStep
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

func drawBlock(block *block) {
	colorDefault := termbox.ColorDefault
	for _, coord := range block.coordinates {
		c := '*'
		termbox.SetCell(coord.x, coord.y, c, colorDefault, colorDefault)
	}
}

func redrawAll(game *game) {
	colorDefault := termbox.ColorDefault
	termbox.Clear(colorDefault, colorDefault)

	drawBox()
	drawBlock(&game.block)

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
		block: block{
			[]coordinate{
				{2, 3},
				{4, 3},
				{6, 3},
				{6, 4},
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
					game.block.moveLeft()
				case termbox.KeyArrowRight:
					game.block.moveRight()
				case termbox.KeyArrowDown:
					game.block.moveDown()
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
