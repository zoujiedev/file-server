package utils

import (
	"encoding/json"
	"io"
	"log"
)

type RespMsg struct {
	Code int         `json:code`
	Msg  string      `json: msg`
	Data interface{} `json:data`
}

func (resp RespMsg) JsonBytes() []byte {
	res, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return res
}

func (resp RespMsg) JsonString() string {
	return string(resp.JsonBytes())
}

func RespSuccess(data interface{}) RespMsg {
	return RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: data,
	}
}

func RespFail(errMsg string) RespMsg {
	return RespMsg{
		Code: -1,
		Msg:  errMsg,
	}
}

func WriteSuccess(w io.Writer, data interface{}) {
	w.Write(RespSuccess(data).JsonBytes())
}

func WriteFail(w io.Writer, msg string) {
	w.Write(RespFail(msg).JsonBytes())
}
