package main

import (
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	
	termbox.SetInputMode(termbox.InputEsc)
	
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()

	x, y := 0, 10
	X, _ := termbox.Size()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				break loop
			}
			if x == X {
				for i:=0; i<X; i++ {
					termbox.SetCell(i, y, ' ', termbox.ColorWhite, termbox.ColorDefault)
				}
				x = 0
			}
			termbox.SetCell(x, y, ev.Ch, termbox.ColorWhite, termbox.ColorDefault)
			termbox.Flush()
			x += 1
		}

	}
}
