package main

import (
	"github.com/nsf/termbox-go"
	"log"
	"os"
)

const (
	tabstop_length            = 8
)

type madman struct {
	termbox_event chan termbox.Event
	buffer *buffer
	quitflag bool
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

func (m *madman) on_key(ev *termbox.Event) {
	switch ev.Key {
	case termbox.KeyCtrlQ:
		m.quitflag = true
	}
}	

func (m *madman) resize() {
}

func (m *madman) handle_event(ev *termbox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		m.on_key(ev)
		if m.quitflag {
			return false
		}
	case termbox.EventResize:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		m.resize()
	case termbox.EventError:
		log.Fatal(ev.Err)
	}
	
	return true
}

func (m *madman) draw() {
}

func new_madman() *madman {
	m := new(madman)
	m.buffer = new_empty_buffer()
	return m
}

func main() {
	logfile, err := os.Create("log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	if err = termbox.Init(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("termbox.Init()")
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
