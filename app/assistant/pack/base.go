package pack

import "github.com/west2-online/DomTok/app/assistant/model"

type _ResponseFactory struct{}

// ResponseFactory is a global variable that points to an instance of _ResponseFactory.
var ResponseFactory _ResponseFactory

func (_ResponseFactory) Error(err error) []byte {
	resp := model.NewResponse()
	resp.SetMeta("type", "error")
	resp.SetData("message", err.Error())
	return resp.MustMarshal()
}

func (_ResponseFactory) Command(command string) []byte {
	resp := model.NewResponse()
	resp.SetMeta("type", "command")
	resp.SetData("command", command)
	return resp.MustMarshal()
}

func (_ResponseFactory) Message(data map[string]interface{}) []byte {
	resp := model.NewResponse()
	resp.SetMeta("type", "message")
	resp.Data = data
	return resp.MustMarshal()
}
