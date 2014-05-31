package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type dynvalFlag int

const (
	dynClientID dynvalFlag = 1 << iota
	dynUUID
	dynRandInt
)

type runner interface {
	run(int, redis.Conn, io.Writer) (int, int)
}

type command struct {
	name string
	args []interface{}

	// Indices for dynamic placeholders
	dynIx map[int]dynvalFlag
}

func (cmd command) run(id int, c redis.Conn, w io.Writer) (int, int) {
	// Insert dynamic values, if any
	args := cmd.args
	if len(cmd.dynIx) > 0 {
		args = make([]interface{}, len(cmd.args))
		for i, arg := range cmd.args {
			if fl, ok := cmd.dynIx[i]; ok {
				sarg := arg.(string)
				if fl&dynClientID == dynClientID {
					sarg = strings.Replace(sarg, "%c", strconv.Itoa(id), -1)
				}
				if fl&dynUUID == dynUUID {
					sarg = strings.Replace(sarg, "%u", uuid.NewRandom().String(), -1)
				}
				if fl&dynRandInt == dynRandInt {
					sarg = strings.Replace(sarg, "%d", strconv.Itoa(rand.Int()), -1)
				}
				args[i] = sarg
			} else {
				args[i] = arg
			}
		}
	}

	// Execute the command
	begin := time.Now()
	res, err := c.Do(cmd.name, args...)
	end := time.Now()
	// TODO: Logging time is interesting, but no more hash to compare results :(

	// Log the results
	if bres, ok := res.([]byte); ok {
		w.Write([]byte(fmt.Sprintf("%d [%s]: %s %v | %#v | %#v\n", id, end.Sub(begin), cmd.name, args, string(bres), err)))
	} else {
		w.Write([]byte(fmt.Sprintf("%d [%s]: %s %v | %#v | %#v\n", id, end.Sub(begin), cmd.name, args, res, err)))
	}

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

func extractCommand(p string, cmd map[string]interface{}) (command, error) {
	var newc command

	if len(cmd) != 1 {
		return newc, fmt.Errorf("%s: invalid JSON file (command must have a single key)", p)
	}
	for k, v := range cmd {
		args, ok := v.([]interface{})
		if !ok {
			return newc, fmt.Errorf("%s: invalid JSON file (command args must be an array)", p)
		}
		newc = command{name: k, args: args, dynIx: make(map[int]dynvalFlag)}
		for i, arg := range args {
			if s, ok := arg.(string); ok {
				var fl dynvalFlag
				if strings.Contains(s, "%c") {
					fl |= dynClientID
				}
				if strings.Contains(s, "%d") {
					fl |= dynRandInt
				}
				if strings.Contains(s, "%u") {
					fl |= dynUUID
				}
				if fl != 0 {
					newc.dynIx[i] = fl
				}
			}
		}
	}
	return newc, nil
}

func extractCommands(p string, cmds []interface{}) ([]runner, int, error) {
	var rs []runner
	var cnt int

	for _, cmd := range cmds {
		switch cmd := cmd.(type) {
		case map[string]interface{}:
			// Extract plain commands
			newc, err := extractCommand(p, cmd)
			if err != nil {
				return nil, 0, err
			}
			rs = append(rs, newc)
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

	// Open the source file
	f, err := os.Open(p)
	if err != nil {
		return file, err
	}
	defer f.Close()

	// Decode the JSON
	var cmds []interface{}
	err = json.NewDecoder(f).Decode(&cmds)
	if err != nil {
		return file, err
	}

	// Extract the commands
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
