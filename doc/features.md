# gred - feature coverage

## Commands

This section lists the implemented commands. `√` means fully supported, `ø` means not supported, and `≈` means partially supported.

The commands are listed by category, as it is on the [redis website][redis].

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
| SETRANGE         | ø      |                                        |
| STRLEN           | √      |                                        |

### Hashes

| Command          | Status | Comment                                |
| ---------------- | :----: | -------------------------------------- |
| HDEL             | ≈      | Does not remove the key once empty.    |
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
| BLPOP            | ø      | |
| BRPOP            | ø      | |
| BRPOPLPUSH       | ø      | |
| LINDEX           | √      | |
| LINSERT          | √      | |
| LLEN             | √      | |
| LPOP             | ≈      | Does not remove key if list is empty.  |
| LPUSH            | ≈      | Seems buggy (order of the elements).   |
| LPUSHX           | √      | |
| LRANGE           | √      | |
| LREM             | ≈      | Does not remove key if list is empty.  |
| LSET             | √      | |
| LTRIM            | √      | |
| RPOP             | √      | |
| RPOPLPUSH        | √      | |
| RPUSH            | √      | |
| RPUSHX           | √      | |

### Sets
