package main

import (
	"io"

	"github.com/garyburd/redigo/redis"
)

type runner interface {
	run(int, redis.Conn, io.Writer) (int, int)
}

type command struct {
	name string
	args []interface{}

	// TODO: Pre-parse the args, prepare for dynamic placeholders
}

func (cmd command) run(id int, c redis.Conn, w io.Writer) (int, int) {
	// Execute the command
	_, err := c.Do(cmd.name, cmd.args...)

	// Log the results
	//writeResults(w, ret, err)

	// Return the number of commands executed, and the number of errors
	if err != nil {
		return 1, 1
	}
	return 1, 0
}

type pipeline struct {
	cmds []command
}

func (p pipeline) run(id int, c redis.Conn, w io.Writer) (int, int) {
	var nerr int

	for _, cmd := range p.cmds {
		err := c.Send(cmd.name, cmd.args...)
		if err != nil {
			//TODO : Log error
			nerr++
		}
	}
	_, err := c.Do("")
	if err != nil {
		nerr++
	}
	return len(p.cmds), nerr
}

type jsonFile struct {
	rs []runner
}

func (j jsonFile) exec(id int, c redis.Conn, w io.Writer, stop <-chan bool) (int, int) {
	ccnt, ecnt := 0, 0
loop:
	for _, r := range j.rs {
		select {
		case <-stop:
			break loop
		default:
		}
		nc, ne := r.run(id, c, w)
		ccnt += nc
		ecnt += ne
	}
	return ccnt, ecnt
}
