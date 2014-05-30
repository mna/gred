package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

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
	res, err := c.Do(cmd.name, cmd.args...)

	// Log the results
	w.Write([]byte(fmt.Sprintf("%d: %s %v | %#v | %#v\n", id, cmd.name, cmd.args, res, err)))

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
			w.Write([]byte(err.Error()))
			w.Write([]byte{'\n'})
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
	rs   []runner
	cmds int
}

// TODO : How to stop a running, maybe blocked, command on a stop signal?

func (j jsonFile) exec(id int, c redis.Conn, w io.Writer, stop <-chan struct{}) (int, int) {
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

func extractCommands(p string, cmds []interface{}) ([]runner, int, error) {
	var rs []runner
	var cnt int

	for _, cmd := range cmds {
		switch cmd := cmd.(type) {
		case map[string]interface{}:
			// Extract plain commands
			if len(cmd) != 1 {
				return nil, 0, fmt.Errorf("%s: invalid JSON file (command must have a single key)", p)
			}
			for k, v := range cmd {
				args, ok := v.([]interface{})
				if !ok {
					return nil, 0, fmt.Errorf("%s: invalid JSON file (command args must be an array)", p)
				}
				rs = append(rs, command{k, args})
			}
			cnt++

		case []interface{}:
			// Extract pipelined commands
			pl := pipeline{}
			list, _, err := extractCommands(p, cmd)
			if err != nil {
				return nil, 0, err
			}
			// Convert []runner to []command
			pl.cmds = make([]command, len(list))
			for i, l := range list {
				c, ok := l.(command)
				if !ok {
					return nil, 0, fmt.Errorf("%s: invalid JSON file (pipeline must contain only commands)", p)
				}
				pl.cmds[i] = c
				cnt++
			}
			rs = append(rs, pl)

		default:
			return nil, 0, fmt.Errorf("%s: invalid JSON file (type %T)", p, cmd)
		}
	}

	return rs, cnt, nil
}

func loadJSONFile(p string) (jsonFile, error) {
	file := jsonFile{}
	f, err := os.Open(p)
	if err != nil {
		return file, err
	}
	defer f.Close()

	var cmds []interface{}
	err = json.NewDecoder(f).Decode(&cmds)
	if err != nil {
		return file, err
	}

	list, cnt, err := extractCommands(p, cmds)
	if err != nil {
		return file, err
	}
	file.rs = list
	file.cmds = cnt
	return file, nil
}

func loadJSONFiles(paths []string) ([]jsonFile, error) {
	files := make([]jsonFile, len(paths))
	for i, p := range paths {
		f, err := loadJSONFile(p)
		if err != nil {
			return nil, err
		}
		files[i] = f
	}
	return files, nil
}
