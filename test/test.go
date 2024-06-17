package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type AuthmRequest struct {
	Op   int
	Args []string
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

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println(v.Subject)
	}
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

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

	reply := make([]byte, 256)
	_, err = conn.Read(reply)
	log.Print("ANSW: " + string(reply))
	log.Print("client: exiting")
}
