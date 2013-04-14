package main

type line struct {
	data []byte
	next *line	
	prev *line
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


