package api

// ResponseMsg : a json response with message.
type ResponseMsg struct {
	Msg string `json:"status"`
}

// ErrorMsg : a json response with message.
type ErrorMsg struct {
	Msg string `json:"error"`
}

// OkResponse is a json ok response.
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
