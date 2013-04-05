package main

import (
	"github.com/nsf/termbox-go"
	"strings"
	"log"
	"os"
)

func draw_n(x, y, count int, color termbox.Attribute, ch rune) {
	for i := 0; i < count; i++ {
		termbox.SetCell(x+i, y, ch, color, termbox.ColorDefault)
	}
}

type Ev struct {
	termbox.Event
	ttype string
}

func keyb() chan Ev {
	res := make(chan Ev)
	var ev termbox.Event
	go func() {
	loop:
		for {
			switch ev = termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyCtrlQ {
					break loop
				}
				ttype := "insert"
				if ev.Key == termbox.KeyBackspace ||
				   ev.Key == termbox.KeyBackspace2 {
					ttype = "backspace"
				} 
				log.Println(ev)
				res <- Ev{ev, ttype}
			}

		}
		res <- Ev{ev, "quit"}
	}()
	return res
}

func draw(s chan Ev) chan bool {
	x, y := 0, 10
	X, Y := termbox.Size()
	cb := termbox.CellBuffer()
	buf := make([]rune, 0, 0)
	res := make(chan bool)
	go func() {
	loop:
		for {
			switch ev := <-s; ev.ttype {
			case "insert":
				if x == X {
					cb = termbox.CellBuffer()
					copy(cb, cb[X:])
					x = 0
				}
				termbox.SetCell(x, y, ev.Ch, termbox.ColorWhite, termbox.ColorBlue)
				if ev.Key == termbox.KeySpace {
					buf = append(buf, ' ')
				} else {
					buf = append(buf, ev.Ch)
				}
				x += 1
				words := strings.Fields(string(buf))
				draw_n(0, Y-1, len(words), termbox.ColorRed, '#')	
				termbox.Flush()
			case "backspace":
				x -= 1
				if x == -1 {
					y = y-1
					x = X-1
				}
				termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorDefault)
				termbox.Flush()
			case "quit":
				break loop

			}

		}
		res <- true
	}()
	return res
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	logfile, err := os.Create("log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
	keychan := keyb()
	<-draw(keychan)

}
