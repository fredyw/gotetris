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
	"sort"
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
	topLeft     coordinate
	topRight    coordinate
	coordinates []coordinate
}

type game struct {
	block block
}

type sortByX []coordinate

func (sbx sortByX) Len() int {
	return len(sbx)
}

func (sbx sortByX) Swap(i, j int) {
	sbx[i], sbx[j] = sbx[j], sbx[i]
}

func (sbx sortByX) Less(i, j int) bool {
	return sbx[i].x < sbx[j].x
}

func (b *block) sortX() {
	sort.Sort(sortByX(b.coordinates))
}

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

func (b *block) rotate() {
	//fmt.Println("before transpose: ", b.coordinates)
	newCoordinates := []coordinate{}
	// transpose the x and y coordinates
	for _, coord := range b.coordinates {
		newX := coord.y
		newY := coord.x
		newCoordinates = append(newCoordinates, coordinate{newX, newY})
	}
	sort.Sort(sortByX(newCoordinates))
	//fmt.Println("after transpose: ", newCoordinates)
	// reverse the x coordinates
	// TODO: incorrect algorithm
	xSize := b.topRight.x - b.topLeft.x
	for i := 0; i < len(newCoordinates); i++ {
		newX := newCoordinates[i].x + xSize
		if newX > b.topRight.x {
			newX = b.topLeft.x + (newX - b.topRight.x - 1)
		}
		newCoordinates[i].x = newX
	}
	//fmt.Println("after reverse:", newCoordinates)
	b.coordinates = newCoordinates
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
	coordMap := map[int]int{}
	i := 0
	block.sortX()
	for _, coord := range block.coordinates {
		c := '*'
		xCoord := coord.x
		if val, ok := coordMap[xCoord]; ok {
			xCoord = val
		} else {
			xCoord += i
			i++
		}
		coordMap[coord.x] = xCoord
		termbox.SetCell(xCoord, coord.y, c, colorDefault, colorDefault)
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
			topLeft:  coordinate{4, 4},
			topRight: coordinate{6, 6},
			coordinates: []coordinate{
				{6, 5},
				{5, 4},
				{6, 4},
				{4, 4},
				//{6, 4},
				//{6, 5},
				//{6, 6},
				//{5, 6},
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
				case termbox.KeySpace:
					game.block.rotate()
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
