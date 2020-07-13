package api

// ResponseMsg : a json response with message
type ResponseMsg struct {
	Msg string `json:"status"`
}

// ErrorMsg : a json response with message
type ErrorMsg struct {
	Msg string `json:"error"`
}

// OkResponse : a json ok reponse
func okResponse(msg ...string) ResponseMsg {
	m := "ok"
	if len(msg) > 0 {
		m = msg[0]
	}
	resp := ResponseMsg{
		Msg: m,
	}
	return resp
}

// ErrorResponse : a json error reponse
func errorResponse(msg string) ErrorMsg {
	resp := ErrorMsg{
		Msg: msg,
	}
	return resp
}
