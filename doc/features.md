# gred - feature coverage

Legend: `√` means fully supported, `ø` means not supported, and `≈` means partially supported (see comment for details).

## High-Level Features

This is a *tl;dr;* version of the 2.8 Redis-compatibility status of the project.

* `redis-cli` and RESP-based clients compatibility: √
* Pipelining: ø
* Telnet: ø
* Clustering, sharding, partitioning, replication, twemproxy support: ø
* Signal handling: ø
* Persistence: ø
* Configuration: ø
* Limits checks (like 512Mb values limit, and offset/indices args): ø

The commands support is detailed in the next section.

## Commands

This section lists the implemented commands. 

The commands are listed by category, as it is on the [redis website][redis].

**TODO** : Must set zero value to elements of a slice before slicing them out (i.e. before `sl = sl[1:]`). Will leak memory otherwise.

### Keys

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| DEL              | √      |                                        |
| DUMP             | ø      |                                        |
| EXISTS           | √      |                                        |
| EXPIRE           | √      |                                        |
| EXPIREAT         | √      |                                        |
| KEYS             | ø      |                                        |
| MIGRATE          | ø      |                                        |
| MOVE             | ø      |                                        |
| OBJECT           | ø      |                                        |
| PERSIST          | √      |                                        |
| PEXPIRE          | √      |                                        |
| PEXPIREAT        | √      |                                        |
| PTTL             | √      |                                        |
| RANDOMKEY        | ø      |                                        |
| RENAME           | ø      |                                        |
| RENAMENX         | ø      |                                        |
| RESTORE          | ø      |                                        |
| SCAN             | ø      |                                        |
| SORT             | ø      |                                        |
| TTL              | √      |                                        |
| TYPE             | √      |                                        |

### Strings

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| APPEND           | √      |                                        |
| BITCOUNT         | ø      |                                        |
| BITOP            | ø      |                                        |
| BITPOS           | ø      |                                        |
| DECR             | √      | Converted to int on each execution.    |
| DECRBY           | √      | Converted to int on each execution.    |
| GET              | √      |                                        |
| GETBIT           | ø      |                                        |
| GETRANGE         | √      |                                        |
| GETSET           | √      |                                        |
| INCR             | √      | Converted to int on each execution.    |
| INCRBY           | √      | Converted to int on each execution.    |
| INCRBYFLOAT      | √      | Converted to float on each execution (like Redis?). |
| MGET             | ø      |                                        |
| MSET             | ø      |                                        |
| MSETNX           | ø      |                                        |
| PSETEX           | ø      |                                        |
| SET              | ≈      | Optional args not implemented (EX, PX, NX, XX). |
| SETBIT           | ø      |                                        |
| SETEX            | ø      |                                        |
| SETNX            | ø      |                                        |
| SETRANGE         | √      |                                        |
| STRLEN           | √      |                                        |

### Hashes

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| HDEL             | √      | Removes the key once empty.            |
| HEXISTS          | √      |                                        |
| HGET             | √      |                                        |
| HGETALL          | √      |                                        |
| HINCRBY          | √      |                                        |
| HINCRBYFLOAT     | √      |                                        |
| HKEYS            | √      |                                        |
| HLEN             | √      |                                        |
| HMGET            | √      |                                        |
| HMSET            | √      |                                        |
| HSCAN            | ø      |                                        |
| HSET             | √      |                                        |
| HSETNX           | √      |                                        |
| HVALS            | √      |                                        |

### Lists

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| BLPOP            | √      | Removes the key once empty.            |
| BRPOP            | √      | Removes the key once empty.            |
| BRPOPLPUSH       | √      | Removes the key once empty.            |
| LINDEX           | √      | |
| LINSERT          | √      | |
| LLEN             | √      | |
| LPOP             | √      | Removes the key once empty.            |
| LPUSH            | √      | |
| LPUSHX           | √      | |
| LRANGE           | √      | |
| LREM             | √      | Removes the key once empty.            |
| LSET             | √      | |
| LTRIM            | √      | Removes the key once empty.            |
| RPOP             | √      | Removes the key once empty.            |
| RPOPLPUSH        | √      | Removes the src key once empty.        |
| RPUSH            | √      | |
| RPUSHX           | √      | |

### Sets

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| SADD             | √      | |
| SCARD            | √      | |
| SDIFF            | √      | |
| SDIFFSTORE       | √      | |
| SINTER           | ø      | |
| SINTERSTORE      | ø      | |
| SISMEMBER        | √      | |
| SMEMBERS         | √      | |
| SMOVE            | ø      | |
| SPOP             | ø      | * |
| SRANDMEMBER      | ø      | |
| SREM             | ø      | * |
| SSCAN            | ø      | |
| SUNION           | ø      | |
| SUNIONSTORE      | ø      | |

### Sorted Sets

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| ZADD             | ø      | |
| ZCARD            | ø      | |
| ZCOUNT           | ø      | |
| ZINCRBY          | ø      | |
| ZINTERSTORE      | ø      | |
| ZLEXCOUNT        | ø      | |
| ZRANGE           | ø      | |
| ZRANGEBYLEX      | ø      | |
| ZRANGEBYSCORE    | ø      | |
| ZRANK            | ø      | |
| ZREM             | ø      | |
| ZREMRANGEBYLEX   | ø      | |
| ZREMRANGEBYRANK  | ø      | |
| ZREMRANGEBYSCORE | ø      | |
| ZREVRANGE        | ø      | |
| ZREVRANGEBYSCORE | ø      | |
| ZREVRANK         | ø      | |
| ZSCAN            | ø      | |
| ZSCORE           | ø      | |
| ZUNIONSTORE      | ø      | |

### HyperLogLog

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| PFADD            | ø      | |
| PFCOUNT          | ø      | |
| PFMERGE          | ø      | |

### Pub/Sub

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| PSUBSCRIBE       | ø      | |
| PUBLISH          | ø      | |
| PUBSUB           | ø      | |
| PUNSUBSCRIBE     | ø      | |
| SUBSCRIBE        | ø      | |
| UNSUBSCRIBE      | ø      | |

### Transactions

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| DISCARD          | ø      | |
| EXEC             | ø      | |
| MULTI            | ø      | |
| UNWATCH          | ø      | |
| WATCH            | ø      | |

### Scripting

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| EVAL             | ø      | |
| EVALSHA          | ø      | |
| SCRIPT EXISTS    | ø      | |
| SCRIPT FLUSH     | ø      | |
| SCRIPT KILL      | ø      | |
| SCRIPT LOAD      | ø      | |

### Connection

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| AUTH             | ø      | |
| ECHO             | √      | |
| PING             | √      | |
| QUIT             | √      | |
| SELECT           | ø      | |

### Server

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| BGREWRITEAOF     | ø      | |
| BGSAVE           | ø      | |
| CLIENT GETNAME   | ø      | |
| CLIENT KILL      | ø      | |
| CLIENT LIST      | ø      | |
| CLIENT PAUSE     | ø      | |
| CLIENT SETNAME   | ø      | |
| CONFIG GET       | ø      | |
| CONFIG RESETSTAT | ø      | |
| CONFIG REWRITE   | ø      | |
| CONFIG SET       | ø      | |
| DBSIZE           | ø      | |
| DEBUG OBJECT     | ø      | |
| DEBUG SEGFAULT   | ø      | |
| FLUSHALL         | ø      | |
| FLUSHDB          | ø      | |
| INFO             | ø      | |
| LASTSAVE         | ø      | |
| MONITOR          | ø      | |
| SAVE             | ø      | |
| SHUTDOWN         | ø      | |
| SLAVEOF          | ø      | |
| SLOWLOG          | ø      | |
| SYNC             | ø      | |
| TIME             | ø      | |

[redis]: http://redis.io/commands
