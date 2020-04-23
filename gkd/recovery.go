package gkd

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr
	// 此处存储的第0个错误是callers，第1个错误是trace,第2个错误是上层的defer func
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(fmt.Sprintf(message + "\nTraceback:"))
	for _, pc := range pcs[:n]{
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
func Recovery() HandlerFunc {
	return func(context *Context) {
		defer func() {
			if err := recover();err !=nil{
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				context.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		context.Next()
	}
}