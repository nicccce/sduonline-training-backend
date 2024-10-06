package app

import "github.com/gin-gonic/gin"

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
type Wrapper struct {
	Ctx *gin.Context
}

func NewWrapper(c *gin.Context) *Wrapper {
	return &Wrapper{Ctx: c}
}

func (w Wrapper) OK() {
	w.Ctx.JSON(200, Result{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
}
func (w Wrapper) Success(data interface{}) {
	w.Ctx.JSON(200, Result{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}
func (w Wrapper) Error(msg string, code ...int) {
	statusCode := 500 // 默认状态码
	if len(code) > 0 {
		statusCode = code[0]
	}

	w.Ctx.JSON(200, Result{
		Code: statusCode,
		Msg:  msg,
		Data: nil,
	})
}
