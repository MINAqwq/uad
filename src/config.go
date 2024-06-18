package main

import (
	"log"

	"gopkg.in/ini.v1"
)

type GlobalConfig struct {
	ipinterface string
	port        string
	tls_cert    string
	tls_key     string
	db_user     string
	db_pass     string
	db_addr     string
	db_port     string
	salt_size   int
	salt_pos    int
}

var global_cfg GlobalConfig

const KEY_CONFIG_SERVER string = "Server"
const KEY_CONFIG_SERVER_INTERFACE string = "INTERFACE"
const KEY_CONFIG_SERVER_PORT string = "PORT"
const KEY_CONFIG_SERVER_TLS_CERT string = "TLS_CERT"
const KEY_CONFIG_SERVER_TLS_KEY string = "TLS_KEY"

const KEY_CONFIG_DATABASE string = "Database"
const KEY_CONFIG_DATABASE_USER string = "USERNAME"
const KEY_CONFIG_DATABASE_PASS string = "PASSWORD"
const KEY_CONFIG_DATABASE_ADDR string = "ADDRESS"
const KEY_CONFIG_DATABASE_PORT string = "PORT"

const KEY_CONFIG_SECURITY string = "Security"
const KEY_CONFIG_SECURITY_SALT_SIZE string = "SALT_SIZE"
const KEY_CONFIG_SECURITY_SALT_POS string = "SALT_POS"

func config_load(path string) {
	var conf GlobalConfig

	log.Println("[CONFIG] Loading...")

	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("[CONFIG] Error: Can't load %s", path)
	}

	if (!cfg.HasSection(KEY_CONFIG_SERVER)) || (!cfg.HasSection(KEY_CONFIG_DATABASE) || (!cfg.HasSection(KEY_CONFIG_SECURITY))) {
		log.Fatalf("[CONFIG] Error: invalid (%s)", path)
	}

	cfg_server := cfg.Section(KEY_CONFIG_SERVER)
	cfg_db := cfg.Section(KEY_CONFIG_DATABASE)
	cfg_sec := cfg.Section(KEY_CONFIG_SECURITY)

	if (!cfg_server.HasKey(KEY_CONFIG_SERVER_INTERFACE)) || (!cfg_server.HasKey(KEY_CONFIG_SERVER_PORT)) ||
		(!cfg_server.HasKey(KEY_CONFIG_SERVER_TLS_CERT)) || (!cfg_server.HasKey(KEY_CONFIG_SERVER_TLS_KEY)) ||
		(!cfg_db.HasKey(KEY_CONFIG_DATABASE_USER)) || (!cfg_db.HasKey(KEY_CONFIG_DATABASE_PASS)) ||
		(!cfg_db.HasKey(KEY_CONFIG_DATABASE_ADDR)) || (!cfg_db.HasKey(KEY_CONFIG_DATABASE_PORT)) ||
		(!cfg_sec.HasKey(KEY_CONFIG_SECURITY_SALT_SIZE)) || (!cfg_sec.HasKey(KEY_CONFIG_SECURITY_SALT_POS)) {
		log.Fatalf("[CONFIG] Config invalid %s", path)
	}

	conf.ipinterface = cfg_server.Key(KEY_CONFIG_SERVER_INTERFACE).MustString("")
	conf.port = cfg_server.Key(KEY_CONFIG_SERVER_PORT).MustString("")
	conf.tls_cert = cfg_server.Key(KEY_CONFIG_SERVER_TLS_CERT).MustString("")
	conf.tls_key = cfg_server.Key(KEY_CONFIG_SERVER_TLS_KEY).MustString("")

	conf.db_user = cfg_db.Key(KEY_CONFIG_DATABASE_USER).MustString("")
	conf.db_pass = cfg_db.Key(KEY_CONFIG_DATABASE_PASS).MustString("")
	conf.db_addr = cfg_db.Key(KEY_CONFIG_DATABASE_ADDR).MustString("")
	conf.db_port = cfg_db.Key(KEY_CONFIG_DATABASE_PORT).MustString("")

	conf.salt_size = cfg_sec.Key(KEY_CONFIG_SECURITY_SALT_SIZE).MustInt(0)
	conf.salt_pos = cfg_sec.Key(KEY_CONFIG_SECURITY_SALT_POS).MustInt(0)

	global_cfg = conf
}
