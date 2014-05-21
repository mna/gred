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
