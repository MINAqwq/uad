package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"log"
)

var global_db *sql.DB

func db_usr_create(username string, email string, passwd string) {
	stmt, err := global_db.Prepare("INSERT INTO usr (username, email, passwd) VALUES (?, ?, ?);")
	if err != nil {
		log.Print("[  DB  ] Failed to prepare statement: db_usr_create")
		return
	}

	_, err = stmt.Exec(username, email, passwd)
	if err != nil {
		log.Print("[  DB  ] Failed to exec statement: db_usr_create")
		return
	}
}

func db_open() {
	connection_string := global_cfg.db_user + `:` + global_cfg.db_pass + "@tcp(" + global_cfg.db_addr + ":" + global_cfg.db_port + ")/mso"
	db, err := sql.Open("mysql", connection_string)
	if err != nil {
		log.Fatalf("failed to open database (%s) [%s]", err, connection_string)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database (%s) [%s]", err, connection_string)
	}

	log.Println("[  DB  ] connected to mso database")

	global_db = db
}
