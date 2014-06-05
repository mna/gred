# gred

`gred` is a pure-Go concurrent implementation of the [Redis][] server. See [features][] for the current state of supported features and commands.

## Installation

Install [go][], then run:

```
$ go get -u github.com/PuerkitoBio/gred/...
```

## Usage

`gred` uses the Redis Serialization Protocol ([RESP][]), so it is a drop-in replacement for Redis. Provided the `$GOPATH/bin` is in your `$PATH`, run:

```
$ gred
```

to start the server on the default port 6379. It uses [glog][] for logging, so the glog flags are available. Type `gred -h` to get the list of options.

Once gred is running, and provided you have a working Redis installation, you can start the redis client to send commands to the server:

```
$ redis-cli
```

Since gred uses the RESP, all Redis clients should be automatically supported (such as [redigo][]).

## dreadis

Under `tools/` is `dreadis`, an automated Redis client. Using JSON command files, this command-line tool can stress-test or validate the correctness of the server. See its [documentation][dreadis] for more details. Some command files exist under the `fixtures/` directory.

## License

The [BSD 3-Clause license][bsd]. See the LICENSE file for details.

[go]: http://golang.org/doc/install
[RESP]: http://redis.io/topics/protocol
[glog]: https://github.com/golang/glog
[Redis]: http://redis.io
[redigo]: https://github.com/garyburd/redigo
[bsd]: http://opensource.org/licenses/BSD-3-Clause
[features]: https://github.com/PuerkitoBio/gred/wiki/Features
[dreadis]: http://godoc.org/github.com/PuerkitoBio/gred/tools/dreadis
