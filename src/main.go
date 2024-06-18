package main

func main() {
	config_load("config.ini")
	db_open()
	server_run()
}
