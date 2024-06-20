package main

func main() {
	config_load("config.ini")
	security_crypt_setup()
	db_open()
	server_run()
}
