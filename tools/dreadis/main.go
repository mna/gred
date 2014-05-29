/*
Command dreadis is an automated tester of RESP-compatible, Redis-like servers.

It supports JSON files that describe the commands to execute, and can run those
scenarios concurrently, with or without pipelining, repeatedly for a given
duration or a given number of iterations.

It saves the replies from the server and generates a hash so that correctness
of replies can easily be validated.

Usage

An example command-line usage is:

    dreadis -c 10 -n 5 -t 30s FILES...

dreadis supports the following flags:

  -c : number of concurrent client connections, defaults to 1.
  -n : number of iterations over the commands in the source files, defaults to 1.
  -t : duration of the test, defaults to the time required for 1 iteration over the files.
  -o : output file to log the replies, defaults to a temporary file, removed on exit.

If -n is not specified and -t is specified, each client will iterate over its source file
for this duration. If both -c and -t are specified, it will stop at the first threshold
that is reached.

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
pipelining. The object should have a single key, the command name, and its
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

An exception to this rule is when the -t flag is used. In this case, execution
of commands may stop at an arbitrary position.

The command generates a sha-256 hash of the replies so that it can be easily
compared, saved and validated. The full log of the replies can be stored in a
file using the -o flag.

Along with the hash, the command produces various statistics, printed to stdout
on completion:

  sha-256: <hash>
  clients: <number of concurrent clients, the -c flag>
  duration: <execution duration>
  iterations: <number of iteration in the files, the -n flag>
  commands: <number of commands executed>
  errors: <number of errors received from the server>
*/
package main
