package main

import (
	"crypto"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"log"
	mrand "math/rand"
)

var salt_symbols = []rune("0123456789abcdef")

var verify_code_symbols = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

var sec_key_pri *rsa.PrivateKey = nil

func security_create_verify_code() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = verify_code_symbols[mrand.Intn(len(verify_code_symbols))]
	}

	return string(b)
}

func security_create_salt() string {
	b := make([]rune, global_cfg.salt_size)
	for i := range b {
		b[i] = salt_symbols[mrand.Intn(len(salt_symbols))]
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
	// log.Printf("[SECRTY] New Hash: %s", hash_str)

	return hash_str
}

func security_hash_extract_salt(pwhash string) string {
	return pwhash[global_cfg.salt_pos:][:global_cfg.salt_size]
}

func security_crypt_setup() {
	priv_key, err := rsa.GenerateKey(crand.Reader, 2024)
	if err != nil {
		log.Fatalln("[SECRTY] " + err.Error())
	}

	// satisfy compiler >:c
	sec_key_pri = priv_key

	log.Println("[SECRTY] RSA Setup")
}

func security_encrypt(content []byte) string {
	crypted, err := rsa.EncryptOAEP(crypto.SHA1.New(), crand.Reader, &sec_key_pri.PublicKey, content, []byte("mso_uad"))
	if err != nil {
		log.Printf("[SECRTY] Failed to encrypt " + string(content))
		return ""
	}

	return hex.EncodeToString(crypted)
}

func security_decrypt(crypted []byte) []byte {
	clear, err := rsa.DecryptOAEP(crypto.SHA1.New(), crand.Reader, sec_key_pri, crypted, []byte("mso_uad"))
	if err != nil {
		log.Printf("[SECRTY] Failed to decrypt " + hex.EncodeToString(crypted))
		return make([]byte, 0)
	}

	return clear
}
