# UAD - User Authentification Daemon

Provides an API reading and returning data in JSON.


## API



### Request Structure

Requests use a predefined object.

```json
{
        Op: (INT)
        Args: [
                "arg1",
                "arg2",
                // ...
        ]
}
```

`op` stands for "Operation" and will tell the Server what request you want to make.


### Response Structure

Responses are also predefined at the first level.

```json
{
        Err: "" // empty string means no error
        Resp: {

        }
}
```

Everything in `resp` (Response) depends on `OP`.


## Operations

`OP` is represented in hexadecimal.

All arguments are strings

| OP   | Name  | Arguments | 
| ---- | ----- | --------- |
| 0000 | VER   | None      |
| 0001 | LOGIN | email, password |
| 0002 | NEW0  | email, username, password |
| 0003 | NEW1  | code |
| 0004 | INFO  | token | 
| 0005 | SAVE  | token, ... |
| 0006 | DEL0  | token |
| 0007 | DEL1  | token , code |

## Database

The UAD speaks with a mysql/mariadb server and reads and writes from a database named `mso`.


### Table

#### usr

| Name | Type |
| ---- | ---- |
| id | `BIGINT UNSIGNED UNIQUE AUTO_INCREMENT NOT NULL` |
| username | `VARCHAR(20) UNIQUE NOT NULL` |
| email | `VARCHAR(40) UNIQUE NOT NULL` |
| passwd | `MEDIUMTEXT UNIQUE NOT NULL` |
| info | `MEDIUMTEXT NOT NULL DEFAULT ""` | 
| created | `TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP` |
| verified | `BOOLEAN NOT NULL DEFAULT FALSE` |

```sql
CREATE TABLE IF NOT EXISTS usr (
        id BIGINT UNSIGNED UNIQUE AUTO_INCREMENT NOT NULL,
        username VARCHAR(20) UNIQUE NOT NULL,
        email VARCHAR(40) UNIQUE NOT NULL,
        passwd MEDIUMTEXT NOT NULL,
        info MEDIUMTEXT NOT NULL DEFAULT "",
        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        verified BOOLEAN NOT NULL DEFAULT FALSE
);
```

```sql
INSERT INTO usr (username, email, passwd) VALUES ("root", "root@minaqwq.dev", "Root123@");
```