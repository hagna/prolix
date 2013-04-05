package main

import (
	"github.com/nsf/termbox-go"
	"log"
	"os"
	"strings"
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

func writePixelFactory(offx, offy int) func(x, y int, ch rune) {
	X, Y := termbox.Size()
	ox := X/2 + offx
	oy := Y/2 + offy
	var res = func(x, y int, ch rune) {
		nx := ox + x
		ny := oy - y
		termbox.SetCell(nx, ny, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
	return res
}

func draw(s chan Ev) chan bool {

	offset := 0
	_, Y := termbox.Size()

	buf := make([]rune, 0, 0)
	res := make(chan bool)
	write_pixel := writePixelFactory(0, 0)
	cville := makeCircleVille(3)
	go func() {
	loop:
		for {
			switch ev := <-s; ev.ttype {
			case "insert":
				x, y, err := cville.getxy(offset)
				if err != nil {
					offset = 0
					x, y, err = cville.getxy(offset)
				}
				write_pixel(x, y, ev.Ch)
				if ev.Key == termbox.KeySpace {
					buf = append(buf, ' ')
				} else {
					buf = append(buf, ev.Ch)
				}
				offset += 1
				words := strings.Fields(string(buf))
				draw_n(0, Y-1, len(words), termbox.ColorRed, '#')
				termbox.Flush()
			case "backspace":
				offset -= 1
				x, y, err := cville.getxy(offset)
				if err != nil {
					offset = 0
					x, y, err = cville.getxy(offset)
				}
				write_pixel(x, y, ' ')
				//				termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorDefault)
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
