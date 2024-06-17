package main

import "log"

type AuthmRequest struct {
	Op   int
	Args []string
}

type AuthmResponse struct {
	Err  string
	Resp map[string]any
}

const (
	OP_VER   = iota
	OP_LOGIN = iota
	OP_NEW   = iota
	OP_INFO  = iota
	OP_SAVE  = iota
	OP_DEL0  = iota
	OP_DEL1  = iota
)

func authm_exec(req *AuthmRequest) AuthmResponse {

	resp := AuthmResponse{}

	log.Printf("READ OP: %d", req.Op)
	switch req.Op {
	case OP_VER:
		break
	case OP_LOGIN:
		break
	case OP_NEW:
		break
	case OP_INFO:
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
