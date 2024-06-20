package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type AuthmRequest struct {
	Op   int
	Args []string
}

type AuthmResponse struct {
	Err  string
	Resp map[string]any
}

func test_ver() {

}

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "127.0.0.1:15100", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}

	defer conn.Close()

	log.Println("client: connected to: ", conn.RemoteAddr())

	req := AuthmRequest{}
	req.Op, err = strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	req.Args = os.Args[2:]

	data, err := json.Marshal(req)
	if err != nil {
		log.Fatalln(err)
	}
	conn.Write(data)

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	resp := AuthmResponse{}
	err = json.Unmarshal(reply[:n], &resp)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Err) != 0 {
		log.Fatal("ERROR: " + resp.Err)
	}

	log.Println(resp.Resp)
}
