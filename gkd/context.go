package gkd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	// 原始对象
	Writer http.ResponseWriter
	Request *http.Request
	// request信息,请求地址请求方法
	Path string
	Method string
	// 路由中的参数
	Params map[string]string
	// write信息
	StatusCode int
}

type H map[string]interface{}

func newContext(w http.ResponseWriter, r *http.Request) *Context{
	return &Context{
		Writer:     w,
		Request:    r,
		Path:       r.URL.Path,
		Method:     r.Method,
	}
}

// 路由从的参数
func (context *Context) Param(key string) string{
	value, _ := context.Params[key]
	return value
}

func (context *Context) Query(key string) string{
	return context.Request.URL.Query().Get(key)
}

// 调用该方法获取表单值
func (context *Context) PostForm(key string) string{
	return context.Request.FormValue(key)
}

// 写入状态码
func (context *Context) Status(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

// 设置请求头
func (context *Context) SetHeader(key string, value string){
	context.Writer.Header().Set(key, value)
}

// 返回string数据
func (context *Context) String (code int, format string, values ...interface{}){
	context.SetHeader("Content-Type", "text/plain")
	context.Status(code)
	// 强制转换传入的值为byte类型，因为write接收byte类型数据
	context.Writer.Write([]byte(fmt.Sprintf(format,values)))
}


// 返回json数据
func (context *Context) JSON (code int, obj interface{}){
	context.SetHeader("Content-Type", "application/json")
	context.Status(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil{
		http.Error(context.Writer, err.Error(), 500)
	}
}

// 返回data
func (context *Context) Data (code int, data []byte){
	context.Status(code)
	context.Writer.Write(data)
}

// 返回html数据
func (context *Context) HTML (code int, html string){
	context.SetHeader("Content-Type", "application/html")
	context.Status(code)
	context.Writer.Write([]byte(html))
}
