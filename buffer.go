package main

import (
	"unicode/utf8"
)

type line struct {
	data []byte
	next *line	
	prev *line
}

// Find a set of closest offsets for a given visual offset
func (l *line) find_closest_offsets(voffset int) (bo, co, vo int) {
	data := l.data
	for len(data) > 0 {
		var vodif int
		r, rlen := utf8.DecodeRune(data)
		data = data[rlen:]
		vodif = rune_advance_len(r, vo)
		if vo+vodif > voffset {
			return
		}

		bo += rlen
		co += 1
		vo += vodif
	}
	return
}
type buffer struct {
	first_line *line
	last_line *line
	lines_n int
	bytes_n int

	words_cache llrb_tree
	words_cache_valid bool

}

func new_empty_buffer() *buffer {
	b := new(buffer)
	l := new(line)
	l.next = nil
	l.prev = nil
	b.first_line = l
	b.last_line = l
	b.lines_n = 1
	return b
}


