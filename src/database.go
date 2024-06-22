package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"log"
)

type UadDbUser struct {
	id            uint64
	username      string
	email         string
	passwd_hashed string
	info          string
	created       string
	verified      bool
}

const (
	DB_QUERY_VERIFY        string = `UPDATE usr SET verified = 1 WHERE id = (SELECT id FROM usr_verify WHERE code = ?);`
	DB_QUERY_DELETE_VERIFY string = `DELETE FROM usr_verify WHERE code = ?;`
	DB_QUERY_USER_EXISTS   string = `SELECT COUNT(id) FROM usr WHERE username = ? OR email = ?;`
	DB_QUERY_CREATE        string = `INSERT INTO usr (username, email, passwd) VALUES (?, ?, ?);`
	DB_QUERY_CREATE_CODE   string = `INSERT INTO usr_verify (id, code) VALUES ((SELECT id FROM usr WHERE username = ?), ?);`
	DB_QUERY_USER_BY_MAIL  string = `SELECT id, username, email, passwd, info, created, verified FROM usr WHERE email = ?;`
	DB_QUERY_USER_BY_ID    string = `SELECT id, username, email, passwd, info, created, verified FROM usr WHERE id = ?;`
	DB_QUERY_UPDATE_INFO   string = `UPDATE usr SET info = ? WHERE id = ?;`
	DB_QUERY_UPDATE_PASSWD string = `UPDATE usr SET passwd = ? WHERE id = ?;`
)

var global_db *sql.DB

func db_usr_get_user_id(id uint64, user *UadDbUser) bool {
	stmt, err := global_db.Prepare(DB_QUERY_USER_BY_ID)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_USER_BY_ID)
		return false
	}

	row := stmt.QueryRow(id)

	err = row.Scan(&user.id, &user.username, &user.email, &user.passwd_hashed, &user.info, &user.created, &user.verified)
	if err != nil {
		log.Printf("[  DB  ] Row scan failed after: %s (%s)", DB_QUERY_USER_BY_ID, err)
		return false
	}

	return true
}

func db_usr_get_user(email string, user *UadDbUser) bool {
	stmt, err := global_db.Prepare(DB_QUERY_USER_BY_MAIL)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_USER_BY_MAIL)
		return false
	}

	row := stmt.QueryRow(email)

	err = row.Scan(&user.id, &user.username, &user.email, &user.passwd_hashed, &user.info, &user.created, &user.verified)
	if err != nil {
		log.Printf("[  DB  ] Row scan failed after: %s (%s)", DB_QUERY_USER_BY_MAIL, err)
		return false
	}

	return true
}

func db_usr_exists(username string, email string) bool {
	stmt, err := global_db.Prepare(DB_QUERY_USER_EXISTS)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_USER_EXISTS)
		return false
	}

	row := stmt.QueryRow(username, email);

	count := 0
	err = row.Scan(&count)
	if err != nil {
		log.Printf("[  DB  ] Row scan failed after: %s (%s)", DB_QUERY_USER_BY_MAIL, err)
		return false
	}

	return (count != 0)

}

// verify a user with the given code
func db_usr_verify(code string) bool {
	stmt, err := global_db.Prepare(DB_QUERY_VERIFY)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_VERIFY)
		return false
	}

	res, err := stmt.Exec(code)
	if err != nil {
		log.Print("[  DB  ] Failed to exec statement: " + DB_QUERY_VERIFY)
		return false
	}

	rows_affected, err := res.RowsAffected()

	if err != nil || rows_affected != 1 {
		return false
	}

	stmt, err = global_db.Prepare(DB_QUERY_DELETE_VERIFY)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_DELETE_VERIFY)
		return false
	}

	_, err = stmt.Exec(code)
	if err != nil {
		log.Print("[  DB  ] Failed to exec statement: " + DB_QUERY_DELETE_VERIFY)
		return false
	}

	return true
}

func db_usr_create_code(username string) bool {
	stmt, err := global_db.Prepare(DB_QUERY_CREATE_CODE)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_CREATE_CODE)
		return false
	}

	for {
		_, err = stmt.Exec(username, security_create_verify_code())
		if err == nil {
			break
		}
		log.Print("[  DB  ] db_usr_create_code exec failed :c (retry...)")
	}

	log.Printf("[  DB  ] Verify code for user %s got created!", username)

	return true
}

// create unverified user account
func db_usr_create(username string, email string, passwd string) bool {
	stmt, err := global_db.Prepare(DB_QUERY_CREATE)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_CREATE)
		return false
	}

	_, err = stmt.Exec(username, email, passwd)
	if err != nil {
		log.Print("[  DB  ] Failed to exec statement: " + DB_QUERY_CREATE)
		return false
	}

	return true
}

func db_usr_update_info(content string, id uint64) bool {
	stmt, err := global_db.Prepare(DB_QUERY_UPDATE_INFO)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_UPDATE_INFO)
		return false
	}

	_, err = stmt.Exec(content, id)
	if err != nil {
		log.Print("[  DB  ] Failed to exec statement: " + DB_QUERY_UPDATE_INFO)
		return false
	}

	return true
}

func db_usr_update_passwd(passwd_hashed string, id uint64) bool {
	stmt, err := global_db.Prepare(DB_QUERY_UPDATE_PASSWD)
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: " + DB_QUERY_UPDATE_PASSWD)
		return false
	}

	_, err = stmt.Exec(passwd_hashed, id)
	if err != nil {
		log.Print("[  DB  ] Failed to exec statement: " + DB_QUERY_UPDATE_PASSWD)
		return false
	}

	return true
}

// connect to database and setup global_db instance
func db_open() {
	connection_string := global_cfg.db_user + `:` + global_cfg.db_pass + "@tcp(" + global_cfg.db_addr + ":" + global_cfg.db_port + ")/mso"
	db, err := sql.Open("mysql", connection_string)
	if err != nil {
		log.Fatalf("Failed to open database (%s) [%s]", err, connection_string)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database (%s) [%s]", err, connection_string)
	}

	log.Println("[  DB  ] Connected to MSO-Database")

	global_db = db
}
