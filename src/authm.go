package main

import (
	"encoding/json"
	"log"
	"time"
)

type AuthmRequest struct {
	Op   int
	Args []string
}

type AuthmResponse struct {
	Err  string
	Resp map[string]any
}

type AuthmUserInfo struct {
	Username string
	Info     string
	Created  string
}

const (
	OP_VER   = iota
	OP_LOGIN = iota
	OP_NEW0  = iota
	OP_NEW1  = iota
	OP_VRFY  = iota
	OP_INFO  = iota
	OP_SAVE  = iota
	OP_DEL0  = iota
	OP_DEL1  = iota
)

func authm_op_ver(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) > 0 {
		resp.Err = "bad arguments"
		return
	}

	resp.Resp["ver"] = 1
}

func authm_op_login(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) != 2 {
		resp.Err = "bad arguments"
		return
	}

	// TODO: check for bad character

	if len(req.Args[0]) < 6 || len(req.Args[0]) > 40 || len(req.Args[1]) < 5 || len(req.Args[1]) > 20 {
		resp.Err = "email or password wrong"
		return
	}

	user := UadDbUser{}
	if !db_usr_get_user(req.Args[0], &user) {
		resp.Err = "email or password wrong"
		return
	}

	passwd_hashed_1 := security_hash_salt_password(req.Args[1], security_hash_extract_salt(user.passwd_hashed))
	if passwd_hashed_1 != user.passwd_hashed {
		resp.Err = "email or password wrong"
		return
	}

	token := session_create(user.id, user.email, int64(time.Hour*2))
	if token == "" {
		resp.Err = "internal error :c"
		return
	}

	resp.Resp["token"] = token
	return
}

func authm_op_new0(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) != 3 {
		resp.Err = "bad arguments"
		return
	}

	// TODO: check for bad character

	// username size check
	if len(req.Args[0]) < 3 || len(req.Args[0]) > 20 {
		resp.Err = "username needs to be between 3 and 20 characters"
		return
	}

	// email size check
	if len(req.Args[1]) < 6 || len(req.Args[1]) > 40 {
		resp.Err = "email needs to be between 6 and 40 characters"
		return
	}

	// password size check
	if len(req.Args[2]) < 5 || len(req.Args[2]) > 20 {
		resp.Err = "password needs to be between 5 and 20 characters"
		return
	}

	if !db_usr_create(req.Args[0], req.Args[1], security_hash_salt_password(req.Args[2], security_create_salt())) {
		resp.Err = "unable to create account (username or email could be taken already)"
		return
	}

	if !db_usr_create_code(req.Args[0]) {
		log.Printf("[ AUTH ] Unable to create verify code for %s <%s>", req.Args[0], req.Args[1])
		resp.Err = "Account created, but we where unable to create a verify code (please contact an admin)"
		return
	}

	log.Printf("[ AUTH ] New Account: %s <%s>", req.Args[0], req.Args[1])
	resp.Resp["msg"] = "Account was created!"
}

func authm_op_new1(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) != 1 {
		resp.Err = "bad arguments"
		return
	}

	if !db_usr_verify(req.Args[0]) {
		resp.Err = "invalid code"
		return
	}

	resp.Resp["msg"] = "Verified ^w^"
}

func authm_op_verify(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) != 1 {
		resp.Err = "bad arguments"
		return
	}

	resp.Resp["valid"] = session_validate(req.Args[0])
}

func authm_op_info(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) != 1 {
		resp.Err = "bad arguments"
		return
	}

	user := UadDbUser{}
	session_data := UserSessionData{}

	if (!session_read(req.Args[0], &session_data)) || (!db_usr_get_user_id(session_data.Id, &user)) {
		resp.Err = "invalid token"
		return
	}

	info := AuthmUserInfo{}
	info.Username = user.username
	info.Info = user.info
	info.Created = user.created

	json_data, err := json.Marshal(info)
	if err != nil {
		resp.Err = "internal error"
		return
	}

	resp.Resp["User"] = json_data
}

func authm_op_save(req *AuthmRequest, resp *AuthmResponse) {
}

func authm_op_del0(req *AuthmRequest, resp *AuthmResponse) {
}

func authm_op_del1(req *AuthmRequest, resp *AuthmResponse) {
}

func authm_exec(req *AuthmRequest) AuthmResponse {

	resp := AuthmResponse{}
	resp.Resp = make(map[string]any)

	log.Printf("[ AUTH ] OP %d", req.Op)

	switch req.Op {
	case OP_VER:
		authm_op_ver(req, &resp)
		break
	case OP_LOGIN:
		authm_op_login(req, &resp)
		break
	case OP_NEW0:
		authm_op_new0(req, &resp)
		break
	case OP_NEW1:
		authm_op_new1(req, &resp)
		break
	case OP_VRFY:
		authm_op_verify(req, &resp)
		break
	case OP_INFO:
		authm_op_info(req, &resp)
		break
	case OP_SAVE:
		break
	case OP_DEL0:
		break
	case OP_DEL1:
		break
	default:
		resp.Err = "invalid operations"
	}

	return resp
}
