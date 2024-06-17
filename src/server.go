package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"net"
)

func _server_handle_client(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("[SERVER] %s | Error: %s", conn.RemoteAddr(), err)
		return
	}

	log.Printf("[SERVER] %s | REQUEST\n%s\n", conn.RemoteAddr(), string(buf))

	req := AuthmRequest{}
	err = json.Unmarshal(buf[:n], &req)
	if err != nil {
		log.Printf("[SERVER] %s, Error: request is not json (%s)", conn.RemoteAddr(), err)
		return
	}

	resp := authm_exec(&req)

	resp_data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[SERVER] %s, Error: marshal response (%s)", conn.RemoteAddr(), err)
		return
	}

	_, err = conn.Write(resp_data)

	if err != nil {
		log.Printf("[SERVER] %s | Error: %s", conn.RemoteAddr(), err)
		return
	}
}

func server_run() {
	log.Printf("[SERVER] KEY : %s", global_cfg.tls_key)
	log.Printf("[SERVER] CERT: %s", global_cfg.tls_cert)

	cert, err := tls.LoadX509KeyPair(global_cfg.tls_cert, global_cfg.tls_key)
	if err != nil {
		log.Fatalf("[SERVER] Error: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	service := global_cfg.ipinterface + ":" + global_cfg.port
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("[SERVER] Error: %s", err)
	}

	log.Print("[SERVER] start listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[SERVER] Error: %s", err)
			break
		}

		defer conn.Close()

		log.Printf("[SERVER] connection from %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}

		go _server_handle_client(conn)
	}
}
