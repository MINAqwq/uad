package main

import (
	"encoding/json"
	"log"
	"strconv"
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
	OP_DEL   = iota
)

func authm_op_ver(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) > 0 {
		resp.Err = "bad arguments"
		return
	}

	resp.Resp["ver"] = 1
}

func authm_op_login(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) != 3 {
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

	if !user.verified {
		resp.Err = "please verify your account"
		return
	}

	passwd_hashed_1 := security_hash_salt_password(req.Args[1], security_hash_extract_salt(user.passwd_hashed))
	if passwd_hashed_1 != user.passwd_hashed {
		resp.Err = "email or password wrong"
		return
	}

	priv_mask, err := strconv.Atoi(req.Args[2])
	if err != nil {
	}

	token := session_create(user.id, user.email, user.passwd_hashed, int64(time.Hour*2), uint32(priv_mask))
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

	if db_usr_exists(req.Args[0], req.Args[1]) ||
		!db_usr_create(req.Args[0], req.Args[1], security_hash_salt_password(req.Args[2], security_create_salt())) {
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

	session_data := UserSessionData{}
	resp.Resp["valid"] = session_read(req.Args[0], &session_data) && session_validate(&session_data)
}

func authm_op_info(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) != 1 {
		resp.Err = "bad arguments"
		return
	}

	user := UadDbUser{}
	session_data := UserSessionData{}

	if (!session_read(req.Args[0], &session_data)) || (!session_validate(&session_data)) {
		resp.Err = "invalid token"
		return
	}

	if (session_data.Privs & SESSION_PRIV_RINFO) == 0 {
		resp.Err = "insufficient permissions"
		return
	}

	db_usr_get_user_id(session_data.Id, &user)

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
	if len(req.Args) != 3 || len(req.Args[2]) == 0 {
		resp.Err = "bad arguments"
		return
	}

	session_data := UserSessionData{}

	if (!session_read(req.Args[0], &session_data)) || (!session_validate(&session_data)) {
		resp.Err = "invalid token"
		return
	}

	if (session_data.Privs & SESSION_PRIV_EXTAC) == 0 {
		resp.Err = "insufficient permissions"
		return
	}

	if req.Args[1] == "info" {
		if db_usr_update_info(req.Args[2], session_data.Id) {
			resp.Resp["msg"] = "done!"
			return
		}
	} else if req.Args[1] == "passwd" {
		if db_usr_update_passwd(security_hash_salt_password(req.Args[2], security_create_salt()), session_data.Id) {
			resp.Resp["msg"] = "done!"
			return
		}
	} else {
		resp.Err = "invalid field"
		return
	}

	resp.Err = "internal error"
	return
}

func authm_op_del(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) != 1 {
		resp.Err = "bad arguments"
		return
	}

	session_data := UserSessionData{}

	if (!session_read(req.Args[0], &session_data)) || (!session_validate(&session_data)) {
		resp.Err = "invalid token"
		return
	}

	if (session_data.Privs & SESSION_PRIV_EXTAC) == 0 {
		resp.Err = "insufficient permissions"
		return
	}

	resp.Resp["msg"] = "Account " + session_data.Email + " got deleted :c"

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
		authm_op_save(req, &resp)
		break
	case OP_DEL:
		authm_op_del(req, &resp)
		break
	default:
		resp.Err = "invalid operations"
	}

	return resp
}
