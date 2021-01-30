package mservice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/o-kit/micro-kit/misc/context"
)

// WebApi 注册
type WebApiRegister interface {
	WebApiRegister(method string, handler func(ctx *context.T, w http.ResponseWriter, req *http.Request))
	WebApiRegisterMethod(service, method string, handler func(ctx *context.T, w http.ResponseWriter, req *http.Request)) bool

	// 解码
	WebApiDecode(ctx *context.T, req *http.Request, arg interface{}) error
	WebApiHandleResp(*context.T, http.ResponseWriter, interface{}, error)
}

// WEB 服务 + 利用http的路由保存method，以供提供WebApi服务
type WebServer struct {
	mux        *http.ServeMux // 路由
	options    Options
	methods    map[string][]string
	middleware []Middleware
}

// 中间件 - 本质就是一个func，将函数保存起来等到需要的时候对入参进行处理
type Middleware func(HandlerFunc) HandlerFunc

func (s *WebServer) AddMiddleware(m ...Middleware) {
	s.middleware = append(s.middleware, m...)
}

type HandlerFunc func(ctx *context.T, w http.ResponseWriter, req *http.Request)

func (s *WebServer) initMux() {
	if s.mux != nil {
		return
	}
	s.mux = http.DefaultServeMux
}

// 提取service名称
func (s *WebServer) getServiceName(method string) string {
	if strings.HasPrefix(method, "/api/") {
		method = method[5:]
	}
	if idx := strings.Index(method, "/"); idx > 0 {
		method = method[:idx]
		return method
	}

	return ""
}

// 提取方法名称
func (s *WebServer) getMethodName(method string) string {
	sp := strings.Split(method, "/")
	return sp[len(sp)-1]
}

// 不重复注册
func (s *WebServer) webapiEnsureMethod(serviceName, methodName string) bool {
	if s.methods == nil {
		s.methods = make(map[string][]string)
	}
	for _, item := range s.methods[serviceName] {
		if methodName == item {
			return false
		}
	}
	s.methods[serviceName] = append(
		s.methods[serviceName], methodName)
	return true
}

// 包装了Handler -> 对于每个http请求 -> 会先执行中间件中的内容
func (s WebServer) wrapHandler(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// 主要是执行之前注册入的中间件
		currentHandleFunc := f
		ctx := context.From(req.Context())
		for _, m := range s.middleware {
			currentHandleFunc = m(currentHandleFunc)
		}
		currentHandleFunc(&ctx, w, req)
	}
}

func (s *WebServer) Ping(w http.ResponseWriter, req *http.Request) {
	hostName, _ := os.Hostname()
	w.Write([]byte(`{"message": ` + strconv.Quote(hostName) + `}`))
}

// 启动服务
func (s *WebServer) Serve(ctx context.T, ln net.Listener) error {
	s.initMux()
	s.mux.HandleFunc("/ping", s.Ping)
	// s.mux.HandleFunc("/debug/pprof/", pprof.Index)
	// s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// s.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	for serviceName := range s.methods {
		s.mux.HandleFunc("/api/"+serviceName+"/ping", s.Ping)
	}
	svr := &http.Server{
		Handler: s.mux,
	}
	return svr.Serve(ln)
}

// 关闭服务
func (s *WebServer) Close() error {
	return nil
}

// 注册接口到路由中 - 具体实现
func (s *WebServer) WebapiHandlerFunc(path string, handler HandlerFunc) {
	s.initMux()
	s.mux.HandleFunc(path, s.wrapHandler(handler))
}

// 注册方法
func (s *WebServer) WebApiRegister(method string, handler func(ctx *context.T, w http.ResponseWriter, req *http.Request)) {
	serviceName := s.getServiceName(method)
	methodName := s.getMethodName(method)
	if s.webapiEnsureMethod(serviceName, methodName) {
		s.WebapiHandlerFunc(method, handler)
	}
}

func (s *WebServer) WebApiRegisterMethod(serviceName, methodName string, handler func(ctx *context.T, w http.ResponseWriter, req *http.Request)) bool {
	// here 可以加一些处理
	path := fmt.Sprintf("/api/%v/%v", serviceName, methodName)
	s.WebapiHandlerFunc(path, handler)

	return true
}

// 下面👇两个方法是涉及到 grpc 方法 《=》 webapi 相互转化问题
func (s *WebServer) WebApiDecode(ctx *context.T, req *http.Request, arg interface{}) error {

	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return err
	}

	if err := s.unmarshalData(arg, body); err != nil {
		return err
	}

	return nil
}

func (s *WebServer) unmarshalData(arg interface{}, body []byte) error {
	if len(body) == 0 {
		return nil
	}
	if body[0] != '{' {
		decoded, _ := base64.URLEncoding.DecodeString(string(body))
		if len(decoded) > 0 {
			body = decoded
		}
	}
	if arg, ok := arg.(proto.Message); ok {
		u := new(jsonpb.Unmarshaler)
		// TODO 这里有个http body 转化为 proto message 问题
		if err := u.Unmarshal(bytes.NewBuffer(body), arg); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(body, arg); err != nil {
			return err
		}
	}
	return nil
}

// 将GRPC结果再转换为HTTP的返回结果
func (s *WebServer) WebApiHandleResp(ctx *context.T, w http.ResponseWriter, resp interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if msg, ok := resp.(proto.Message); ok {
		(&jsonpb.Marshaler{EmitDefaults: true}).Marshal(w, msg)
	} else {
		ret, _ := json.Marshal(resp)
		w.Write(ret)
	}
	return
}
