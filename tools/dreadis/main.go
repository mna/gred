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
       Defaults to no output.

  -no-stats : do not display the execution statistics.
  -net      : network type, defaults to "tcp".
  -addr     : network address, defaults to ":6379".

If -n is 0 (its default) and -t is specified, each client will iterate over its source file
for this duration. If both -n and -t are specified, it will stop at the first threshold
that is reached. If only -n is specified, it will run this number of iterations for each
source file, without time limit.

FILES are JSON command files. Any number of source files may be specified, and they
will be spread across the -c concurrent clients. The logic is simply to assign the
first file to the first client, second file to the second client, etc. until all clients
have at least one file, and all files are executed by at least one client.

So if -c = 1 and there is more than one file specified, the same client will execute
all files. Conversely, if there are many clients and just one file, all clients will
run the same file.

JSON command files

The JSON command files have the following format:

    [
        {"command": ["arg1", "arg2", 2, 3.45]},
        [
            {"pipelined_cmd1": ["arg1", "arg2"]},
            {"pipelined_cmd2": ["arg1"]}
        ],
    ]

Objects on the top-level array are straight commands, executed without
pipelining. The object must have a single key, the command name, and its
value is an array of arguments to pass to the command.

Array values on the top-level array hold the same kind of object values,
each representing a command to execute, but those commands are pipelined.

Special placeholders %d and %q can be used in the string arguments. Those
placeholders have the following meaning:

    %c : replaced with the client id.
    %u : replaced with a random UUID, newly generated on each execution of the command.
    %d : a random integer.

Results

Each client stores the replies from the server, and since commands and clients
are created and executed in a well-defined order, results can be deterministic
and so can be validated against a reference implementation (i.e. against
the official Redis server).

There are caveats to be aware of when comparing results. When the -t flag is set,
commands may stop at an arbitrary position. When dynamic placeholders %u or %d
are used, the results are very likely to be different from one execution to another.
Also, some commands such as TIME will return different results at different times,
by definition.

The command produces various statistics, printed to stdout on completion unless
the -no-stats flag is set:

  clients: <number of concurrent clients, the -c flag>
  duration: <actual execution duration, which may be different than the -t flag>
  iterations: <number of iteration in the files, the -n flag>
  commands: <number of commands executed>
  errors: <number of errors received from the server>
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

var (
	concurrent = flag.Int("c", 1, "number of concurrent clients")
	iterations = flag.Int("n", 0, "number of iterations over the command files")
	timeout    = flag.Duration("t", 0, "timeout or duration of the execution")
	output     = flag.String("o", "", "output file to log server replies")
	format     = flag.String("f", "", "format of the server replies")

	nostats = flag.Bool("no-stats", false, "do not display execution statistics")
	net     = flag.String("net", "tcp", "network type")
	addr    = flag.String("addr", ":6379", "network address")
)

var usage = `Usage: dreadis [options...] FILES...

Options:
  -c  Number of concurrent clients to run. Defaults to 1.
  -f  Format of server replies output, in Go's text/template
      syntax. Defaults to no output of the server replies.
      Predefined formats can be used by name, they are
      'short', 'std' and 'stdtime'.
  -n  Number of iterations over the command files. Defaults
      to 0, which behaves like -n=1 unless -t is set, in which
      case it will execute the commands until -t is reached.
  -o  Output file to log server replies. Defaults to stdout,
      if an output format is specified using -f.
  -t  Maximum duration of the execution if -n is > 0, or
      duration of the execution if -n is set to 0.

  -no-stats Do not display execution statistics.
  -net      Network interface to use to connect to the server.
            Defaults to "tcp".
  -addr     Network address to use to connect to the server.
            Defaults to ":6379".

  The following struct fields are available for the format (-f)
  template, which must use a {{range .}}...{{end}} action since
  a slice of such struct is passed as template value:

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
		"std":     "{{range .}}[{{.ClientID}}] {{.Command}} {{.Args}} | {{.Result}} | {{.Err}}\n{{end}}",
		"stdtime": "{{range .}}[{{.ClientID}}] [{{.Time}}] {{.Command}} {{.Args}} | {{.Result}} | {{.Err}}\n{{end}}",
		"short":   "{{range .}}[{{.ClientID}}] {{.Command}} | {{.Err}}\n{{end}}",
	}
)

func main() {
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
			fmt.Fprint(os.Stderr, err)
			fmt.Fprintln(os.Stderr)
			os.Exit(1)
		}
	}

	// Load command files
	files, err := loadJSONFiles(flag.Args())
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(2)
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
		fmt.Fprint(os.Stderr, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(3)
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
		var out io.Writer
		out = os.Stdout
		if *output != "" {
			f, err := os.Create(*output)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				fmt.Fprintln(os.Stderr)
				os.Exit(4)
			}
			defer f.Close()
			out = f
		}

		res := collectReplies(ws)
		err = tpl.Execute(out, res)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			fmt.Fprintln(os.Stderr)
			os.Exit(5)
		}
		fmt.Fprint(out, "\n")
	}

	// Display stats
	if !*nostats {
		cmds, errs := collectStats(ws)

		// Display the results
		fmt.Printf(`clients:    %d
iterations: %d
duration:   %s
commands:   %d
errors:     %d
`, *concurrent, *iterations, end.Sub(begin), cmds, errs)
	}
}

func printUsage(msg string, exitCode int) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
		fmt.Fprintln(os.Stderr, "")
	}
	flag.Usage()
	os.Exit(exitCode)
}
