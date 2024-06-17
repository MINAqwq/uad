package main

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

func authm_op_ver(req *AuthmRequest, resp *AuthmResponse) {
	if len(req.Args) > 0 {
		resp.Err = "invalid arguments"
		return
	}

	resp.Resp["ver"] = 1
}

func authm_exec(req *AuthmRequest) AuthmResponse {

	resp := AuthmResponse{}
	resp.Resp = make(map[string]any)

	switch req.Op {
	case OP_VER:
		authm_op_ver(req, &resp)
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
