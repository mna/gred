/*
Command dreadis is an automated tester of RESP-compatible, Redis-like servers.

It supports JSON files that describe the commands to execute, and can run those
scenarios concurrently, with or without pipelining, repeatedly for a given
duration or a given number of iterations.

It can optionnally print the replies from the server, so that they can be consulted,
redirected to a file or piped to a hashing command for quick verification of the
correctness of the results.

Usage

An example command-line usage is:

    dreadis -c 10 -n 5 -t 30s FILES...

dreadis supports the following flags:

  -c : number of concurrent client connections, defaults to 1.
  -n : number of iterations over the commands in the source files, defaults to 0, which
       means run each command once, or until -t is reached, if a timeout is set.
  -t : maximum duration of the test, defaults to 0 (no time limit).
  -o : output file to log the replies, defaults to stdout.
  -f : format of the output, using the syntax of the Go package text/template.
       Defaults to the predefined stats output.

  -net   : network type, defaults to "tcp".
  -addr  : network address, defaults to ":6379".

If -n is 0 (its default) and -t is specified, each client will iterate over its source file
for this duration. If both -n and -t are specified, it will stop at the first threshold
that is reached. If only -n is specified, it will run this number of iterations for each
source file, without time limit.

FILES are JSON command files. Any number of source files may be specified, and they
will be spread across the -c concurrent clients. The logic is simply to assign the
first file to the first client, second file to the second client, etc. until all clients
have at least one file, and all files are executed by at least one client.

So if -c = 1 and there is more than one file specified, the same client will execute
all files in sequence. Conversely, if there are many clients and just one file,
all clients will run the same file.

JSON command files

The JSON command files have the following format:

    [
        {"command": ["arg1", "arg2", 2, 3.45]},
        [
            {"pipelined_cmd1": ["arg1", "arg2"]},
            {"pipelined_cmd2": ["arg1"]}
        ]
    ]

Objects on the top-level array are straight commands, executed without
pipelining. The object must have a single key, the command name, and its
value is an array of arguments to pass to the command.

Array values on the top-level array hold the same kind of object values,
each representing a command to execute, but those commands are pipelined.

Special placeholders can be used in the string arguments. Those
placeholders have the following meaning:

    %c : replaced with the client id.
    %u : replaced with a random UUID, newly generated on each execution of the command.
    %d : a random integer, newly generated on each execution of the command.

Results

Each client stores the replies from the server, and since clients and commands
within a client are created and executed in a well-defined order, results can
be deterministic and so can be validated against a reference implementation
(i.e. against the official Redis server).

There are caveats to be aware of when comparing results. When the -t flag is set,
commands may stop at an arbitrary position. When dynamic placeholders %u or %d
are used, the results are very likely to be different from one execution to another.
Also, some commands such as TIME will return different results at different times,
by definition.

The command produces various statistics, available for the -f output format
and printed by the stats predefined format:

  Clients: <number of concurrent clients, the -c flag>
  Duration: <actual execution duration, which may be different than the -t flag>
  Iterations: <number of iteration in the files, the -n flag>
  Commands: <number of commands executed>
  Errors: <number of errors received from the server>
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"text/template"
	"time"
)

// The command-line flags
var (
	concurrent = flag.Int("c", 1, "number of concurrent clients")
	iterations = flag.Int("n", 0, "number of iterations over the command files")
	timeout    = flag.Duration("t", 0, "timeout or duration of the execution")
	output     = flag.String("o", "", "output file to log server replies")
	format     = flag.String("f", "stats", "format of the server replies")

	net  = flag.String("net", "tcp", "network type")
	addr = flag.String("addr", ":6379", "network address")
)

// Command-line usage of the tool
var usage = `Usage: dreadis [options...] FILES...

Options:
  -c  Number of concurrent clients to run. Defaults to 1.
  -f  Format of the printed output, in Go's text/template
      syntax. Defaults to the predefined 'stats' format.
      Predefined formats can be used by name, they are
      'stats', 'err', 'std' and 'stdtime'.
  -n  Number of iterations over the command files. Defaults
      to 0, which behaves like -n=1 unless -t is set, in which
      case it will execute the commands until -t is reached.
  -o  Output file to log server replies. Defaults to stdout,
      if an output format is specified using -f.
  -t  Maximum duration of the execution if -n is > 0, or
      duration of the execution if -n is set to 0.

  -net      Network interface to use to connect to the server.
            Defaults to "tcp".
  -addr     Network address to use to connect to the server.
            Defaults to ":6379".

  The following struct fields are available for the format (-f)
  template, which must use a {{range .Results}}...{{end}} action to
  print 'cmdResult's (an execInfo value is passed to the template):

    type execInfo struct {
    	Clients    int
    	Iterations int
    	Timeout    time.Duration
    	Time       time.Duration
    	Commands   int
    	Errors     int
    	Files      []string
    	Results    []*cmdResult
    }

    type cmdResult struct {
    	ClientID     int
    	Command      string
    	Args         []interface{}
    	Result       interface{}
    	Err          error
    	Time         time.Duration
    	Pipelined    bool
    	PipelineExec bool
    }
`

var (
	fmts = map[string]string{
		"stats":   "Clients: {{.Clients}}\nIterations: {{.Iterations}}\nDuration: {{.Time}}\nCommands: {{.Commands}}\nErrors: {{.Errors}}\n",
		"std":     "{{range .Results}}[{{.ClientID}}] {{.Command}} {{.Args}} | {{.Result}} | {{.Err}}\n{{end}}",
		"stdtime": "{{range .Results}}[{{.ClientID}}] [{{.Time}}] {{.Command}} {{.Args}} | {{.Result}} | {{.Err}}\n{{end}}",
		"err":     "{{range .Results}}[{{.ClientID}}] {{.Command}} | {{.Err}}\n{{end}}",
	}
)

func main() {
	parseFlags()

	// Parse output template, so that it fails before execution
	// in case of problem.
	var tpl *template.Template
	var err error
	if *format != "" {
		if v, ok := fmts[*format]; ok {
			*format = v
		}
		tpl, err = template.New("main").Parse(*format)
		if err != nil {
			printMessage(err.Error(), 1)
		}
	}

	// Load command files
	files, err := loadJSONFiles(flag.Args())
	if err != nil {
		printMessage(err.Error(), 2)
	}

	// Launch workers, dispatching the command files
	start, stop := make(chan struct{}), make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(*concurrent)
	ws, err := launchWorkers(args{
		net:   *net,
		addr:  *addr,
		c:     *concurrent,
		n:     *iterations,
		wg:    &wg,
		start: start,
		stop:  stop,
		files: files})

	if err != nil {
		printMessage(err.Error(), 3)
	}

	// Start the workers
	begin := time.Now()
	close(start)

	// Stop processing after the timeout duration
	if *timeout > 0 {
		go func() {
			<-time.After(*timeout)
			close(stop)
		}()
	}

	// Wait for workers
	wg.Wait()
	end := time.Now()

	// Print output if a format is specified
	if *format != "" {
		res, ncmd, nerr := collectReplies(ws)
		ex := &execInfo{
			Clients:    *concurrent,
			Iterations: *iterations,
			Timeout:    *timeout,
			Time:       end.Sub(begin),
			Commands:   ncmd,
			Errors:     nerr,
			Files:      flag.Args(),
			Results:    res,
		}
		printResults(ex, tpl, *output)
	}
}

func printResults(ex *execInfo, tpl *template.Template, outFile string) {
	var out io.Writer
	out = os.Stdout
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			printMessage(err.Error(), 4)
		}
		defer f.Close()
		out = f
	}

	err := tpl.Execute(out, ex)
	if err != nil {
		printMessage(err.Error(), 5)
	}
}

func parseFlags() {
	// Parse and validate args
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()
	if flag.NArg() < 1 {
		printUsage("", 1)
	}
	if *concurrent < 1 {
		*concurrent = 1
	}
	// -n = 1 if no timeout is set.
	if *timeout == 0 && *iterations == 0 {
		*iterations = 1
	}
}

func printUsage(msg string, exitCode int) {
	if msg != "" {
		printMessage(msg, 0)
	}
	flag.Usage()
	os.Exit(exitCode)
}

func printMessage(msg string, exitCode int) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, "")
	if exitCode > 0 {
		os.Exit(exitCode)
	}
}
