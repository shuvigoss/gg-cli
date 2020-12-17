package common

import "net/http"

const (
	ServerError   = 500 //服务端异常
	ServerSuccess = 200 //服务正常
	BizError      = 400 //业务异常
)

type Result struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}


func IsSuccess(result *Result) bool {
	return result.Status == http.StatusOK
}
