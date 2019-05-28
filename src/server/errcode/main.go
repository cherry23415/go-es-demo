package errcode

import (
	"errors"
	"go-es-demo/src/server/entity"
)

var (
	SQL_EMPTY_DATA = errors.New("数据库不存在数据")
)

//返回结果错误码定义
var (
	SUCCESS      = &ErrCode{0, "ok"}
	SYSTEM_ERROR = &ErrCode{-1, "系统异常"}
	PARAM_ERROR  = &ErrCode{1001, "参数错误"}
)

type ErrCode struct {
	Code int
	Msg  string
}

func (e *ErrCode) Result() *entity.Result {
	return &entity.Result{
		Status: e.Code,
		Msg:    e.Msg,
		Data:   "",
	}
}

func (e *ErrCode) ResultWithData(data interface{}) *entity.Result {
	return &entity.Result{
		Status: e.Code,
		Msg:    e.Msg,
		Data:   data,
	}
}

func (e *ErrCode) ResultWithMsg(msg string) *entity.Result {
	return &entity.Result{
		Status: e.Code,
		Msg:    msg,
		Data:   "",
	}
}

func (e *ErrCode) ReplaceMsg(msg string) *ErrCode {
	return &ErrCode{
		Code: e.Code,
		Msg:  msg,
	}
}
