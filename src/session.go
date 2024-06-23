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
	Passhsh string
	Expires int64
}

func session_create(id uint64, email string, password_hash string, valid_for int64) string {
	session_data := UserSessionData{}
	session_data.Id = id
	session_data.Email = email
	session_data.Passhsh = session_passhsh_slice(password_hash)
	session_data.Expires = (time.Now().Unix() + valid_for)

	session_json, err := json.Marshal(session_data)
	if err != nil {
		log.Printf("[USESSN] Failed to encode json for map containing [id: %d, email: %s, expires: %d]",
			id, email, session_data.Expires)
		return ""
	}

	crypted := security_encrypt(session_json)

	log.Println("[USESSN] Created Session for " + email)

	return crypted
}

func session_passhsh_slice(password_hash string) string {
	return password_hash[(len(password_hash) - 6):]
}

func session_read(token string, buffer *UserSessionData) bool {
	token_bytes, err := hex.DecodeString(token)
	if err != nil {
		log.Println("[USESSN] " + err.Error())
		return false
	}

	clear_bytes := security_decrypt(token_bytes)
	if len(clear_bytes) == 0 {
		return false
	}

	err = json.Unmarshal(clear_bytes, &buffer)
	if err != nil {
		log.Println("[USESSN] " + err.Error())
		return false
	}

	log.Printf("[USESSN] Session decoded: {%d, %s, %s, %d}",
			buffer.Id, buffer.Email, buffer.Passhsh, buffer.Expires)

	return true
}

func session_validate(token string) bool {

	session_data := UserSessionData{}
	db_user := UadDbUser{}

	return session_read(token, &session_data) &&
		db_usr_get_user_id(session_data.Id, &db_user) &&
		(session_data.Expires > time.Now().Unix()) &&
		(session_data.Email == db_user.email) &&
		(session_data.Passhsh == session_passhsh_slice(db_user.passwd_hashed))
}
