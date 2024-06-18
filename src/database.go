package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"log"
)

var global_db *sql.DB

func db_open() {
	connection_string := global_cfg.db_user + `:` + global_cfg.db_pass + "@tcp(" + global_cfg.db_addr + ":" + global_cfg.db_port + ")/mso"
	db, err := sql.Open("mysql", connection_string)
	if err != nil {
		log.Fatalf("failed to open database (%s) [%s]", err, connection_string)
	}

	global_db = db
}
