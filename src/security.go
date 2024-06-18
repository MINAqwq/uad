package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/rand"
)

var salt_symbols = []rune("0123456789abcdef")

func security_create_salt() string {
	b := make([]rune, global_cfg.salt_size)
	for i := range b {
		b[i] = salt_symbols[rand.Intn(len(salt_symbols))]
	}
	return string(b)
}

func security_hash_salt_password(passwd_clear string, salt_str string) string {
	// hash
	sha := sha256.New()
	sha.Write([]byte(passwd_clear))
	sha.Write([]byte(salt_str))

	hash_str := hex.EncodeToString(sha.Sum(nil))

	// insert salt
	hash_str = hash_str[:global_cfg.salt_pos] + salt_str + hash_str[global_cfg.salt_pos:]
	log.Printf("[SECRTY] New Hash: %s", hash_str)

	return hash_str
}
