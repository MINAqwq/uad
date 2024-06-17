# UAD - User Authentification Daemon

Provides an API returning data in CBOR.


## API



### Request Structure

Requests use a predefined object.

```json
{
        op: (INT)
        args: {
                // ...
        }
}
```

`op` stands for "Operation" and will tell the Server what request you want to make.


### Response Structure

Responses are also predefined at the first level.

```json
{
        err: "" // empty string means no error
        resp: {

        }
}
```

Everything in `resp` (Response) depends on `OP`.


## Operations

`OP` is represented in hexadecimal.

| OP   | Name  | Arguments | 
| ---- | ----- | --------- |
| 0000 | VER   | None      |
| 0001 | LOGIN | email: `str`, password: `str` |
| 0002 | NEW   | email: `str`, username: `str`, password: `str` |
| 0003 | INFO  | token: `str` |
| 0004 | SAVE  | token: `str`, ... |
| 0005 | DEL0  | token: `str` |
| 0006 | DEL1  | token: `str`, code: `str` |

