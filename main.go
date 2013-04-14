package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

type madman struct {
	termbox_event chan termbox.Event
}


func (m *madman) main_loop() {
	m.termbox_event = make(chan termbox.Event, 20)
	go func() {
		for {
			m.termbox_event <- termbox.PollEvent()
		}
	}()
	for {
		select {
		case ev := <-m.termbox_event:
			ok := m.handle_event(&ev)
			if !ok {
				return
			}
			m.consume_more_events()
			m.draw()
			termbox.Flush()
		}
	}
}

func (m *madman) consume_more_events() {
loop:
	for {
		select {
		case ev := <-m.termbox_event:
			ok := m.handle_event(&ev)
			if !ok {
				break loop	
			}
		default:
			break loop
		}
	}
}	

func (m *madman) handle_event(ev *termbox.Event) bool {
	return false
}

func (m *madman) draw() {
}

func new_madman() *madman {
	return new(madman)
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt)

	madman := new_madman()
	//madman.resize()
	madman.draw()
	termbox.SetCursor(0,0)
	termbox.Flush()
	madman.main_loop()
}
