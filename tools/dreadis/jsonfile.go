package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
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
	run(int, redis.Conn) []*cmdResult
}

type command struct {
	name string
	args []interface{}

	// Indices for dynamic placeholders
	dynIx map[int]dynvalFlag
}

func (cmd command) prepareArgs(id int) []interface{} {
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
	return args
}

func (cmd command) run(id int, c redis.Conn) []*cmdResult {
	args := cmd.prepareArgs(id)
	// Execute the command
	begin := time.Now()
	res, err := c.Do(cmd.name, args...)
	end := time.Now()

	cr := &cmdResult{
		ClientID: id,
		Command:  cmd.name,
		Args:     args,
		Result:   cleanResults(res),
		Err:      err,
		Time:     end.Sub(begin),
	}

	return []*cmdResult{cr}
}

type pipeline struct {
	cmds []command
}

func (p pipeline) run(id int, c redis.Conn) []*cmdResult {
	crs := make([]*cmdResult, len(p.cmds)+1)
	for i, cmd := range p.cmds {
		args := cmd.prepareArgs(id)
		err := c.Send(cmd.name, args...)
		crs[i] = &cmdResult{
			ClientID:  id,
			Command:   cmd.name,
			Args:      args,
			Err:       err,
			Pipelined: true,
		}
	}

	// Execute all pipelined commands
	begin := time.Now()
	res, err := redis.Values(c.Do(""))
	end := time.Now()

	// Store the pipeline exec results
	crs[len(crs)-1] = &cmdResult{
		ClientID:     id,
		Err:          err,
		PipelineExec: true,
		Time:         end.Sub(begin),
	}

	if err == nil {
		// Save results for each command
		for i, r := range res {
			crs[i].Result = cleanResults(r)
		}
	}

	return crs
}

func cleanResults(res interface{}) interface{} {
	switch res := res.(type) {
	case []byte:
		return string(res)
	case []interface{}:
		for i := 0; i < len(res); i++ {
			res[i] = cleanResults(res[i])
		}
		return res
	default:
		return res
	}
}

type jsonFile struct {
	rs   []runner
	cmds int
}

// TODO : How to stop a running, maybe blocked, command on a stop signal?

// exec executes all commands in this jsonFile, stopping once the stop channel
// is signaled. It returns the number of commands executed, and the number of
// errors returned from the server.
func (j jsonFile) exec(id int, c redis.Conn, stop <-chan struct{}) []*cmdResult {
	var res []*cmdResult
loop:
	for _, r := range j.rs {
		select {
		case <-stop:
			break loop
		default:
		}
		crs := r.run(id, c)
		res = append(res, crs...)
	}
	return res
}

func extractCommand(p string, cmd map[string]interface{}) (command, error) {
	var newc command

	// Must have a single key - the command name
	if len(cmd) != 1 {
		return newc, fmt.Errorf("%s: invalid JSON file (command must have a single key)", p)
	}

	// Loop through the map, it only has a single key
	for k, v := range cmd {
		args, ok := v.([]interface{})
		if !ok {
			return newc, fmt.Errorf("%s: invalid JSON file (command args must be an array)", p)
		}

		// Create the command and inspect arguments for dynamic values placeholders
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
		break
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
