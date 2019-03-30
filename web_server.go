package go_webkit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context 类型别名
type Context = *gin.Context

// WebServer web服务辅助工具
type WebServer struct {
	host   string
	engine *gin.Engine
}

// NewWebServerDefault 构建web服务数据对象
// domain可以为空
func NewWebServerDefault(domain string, port string) *WebServer {
	return &WebServer{fmt.Sprintf("%s:%s", domain, port), gin.Default()}
}

// Run 运行web服务
func (sf *WebServer) Run() error {
	return sf.engine.Run(sf.host)
}

// HandleFunc 监听函数
func (sf *WebServer) HandleFunc(method Method, relativePath string, handle func(ctx Context)) *WebServer {
	switch method {
	case Get:
		sf.engine.GET(relativePath, handle)
	case Post:
		sf.engine.POST(relativePath, handle)
	case Put:
		sf.engine.PUT(relativePath, handle)
	case Delete:
		sf.engine.DELETE(relativePath, handle)
	case Patch:
		sf.engine.PATCH(relativePath, handle)
	case Options:
		sf.engine.OPTIONS(relativePath, handle)
	}
	return sf
}

// HandleFuncGet 监听Get
func (sf *WebServer) HandleFuncGet(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Get, relativePath, handle)
}

// HandleFuncPost 监听Post
func (sf *WebServer) HandleFuncPost(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Post, relativePath, handle)
}

// HandleFuncPut 监听Put
func (sf *WebServer) HandleFuncPut(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Put, relativePath, handle)
}

// HandleFuncDelete 监听Delete
func (sf *WebServer) HandleFuncDelete(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Delete, relativePath, handle)
}

// HandleFuncPatch 监听Patch
func (sf *WebServer) HandleFuncPatch(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Patch, relativePath, handle)
}

// HandleFuncOptions 监听Options
func (sf *WebServer) HandleFuncOptions(relativePath string, handle func(ctx Context)) *WebServer {
	return sf.HandleFunc(Options, relativePath, handle)
}

// HandleStruct 监听结构体，反射街头的http方法以及遍历每个字段的http方法，实现REST形式的API服务
// 结构体的方法必须与 Method 类型的名称一致
func (sf *WebServer) HandleStruct(relativePath string, handle interface{}) *WebServer {
	sf.handleStruct(relativePath, handle)
	return sf
}

// handleStruct 使用反射遍历结构体的方法和字段，对http方法进行注册
func (sf *WebServer) handleStruct(relativePath string, handle interface{}) {
	rfType := reflect.TypeOf(handle)
	rfValue := reflect.ValueOf(handle)
	// 只接受结构体、接口及指针
	switch rfType.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Struct:
		// 反射方法
		for i := 0; i < rfType.NumMethod(); i++ {
			methodName := rfType.Method(i).Name
			if !isMethod(NewMethod(methodName)) {
				continue
			}
			handleFunc := func(ctx Context) {
				rfValue.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(ctx)})
			}
			sf.HandleFunc(NewMethod(methodName), relativePath, handleFunc)
		}
		// 反射字段
		for i := 0; i < rfType.NumField(); i++ {
			fieldName := rfType.Field(i).Name
			sf.handleStruct(relativePath+"/"+strings.ToLower(fieldName[:1])+fieldName[1:],
				rfValue.FieldByName(fieldName).Interface())
		}
	}
}
