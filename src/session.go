package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
)

type UserSessionData struct {
	Id      uint64
	Email   string
	Expires int64
}

func session_create(id uint64, email string, valid_for int64) string {
	session_data := UserSessionData{}
	session_data.Id = id
	session_data.Email = email
	session_data.Expires = (time.Now().Unix() + valid_for)

	session_json, err := json.Marshal(session_data)
	if err != nil {
		log.Printf("[USESSN] Failed to encode json for map containing [id: %d, email: %s, expires: %d]", id, email, session_data.Expires)
		return ""
	}

	crypted := security_encrypt(session_json)

	log.Println("[USESSN] Created Session for " + email)

	return crypted
}

func session_validate(token string) bool {
	token_bytes, err := hex.DecodeString(token)
	if err != nil {
		log.Println("[USESSN] " + err.Error())
		return false
	}

	clear_bytes := security_decrypt(token_bytes)
	if len(clear_bytes) == 0 {
		return false
	}

	session_data := UserSessionData{}
	err = json.Unmarshal(clear_bytes, &session_data)
	if err != nil {
		log.Println("[USESSN] " + err.Error())
		return false
	}

	log.Printf("[USESSN Session decoded: {%d, %s, %d}]", session_data.Id, session_data.Email, session_data.Expires)

	return session_data.Expires > time.Now().Unix()
}
