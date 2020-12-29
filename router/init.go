package router

import (
	"otter-cloud-ws/acl"
	"otter-cloud-ws/config"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/service/apihandler"

	"github.com/EricChiou/httprouter"
	"github.com/valyala/fasthttp"
)

// Init init api
func Init() {
	initUserAPI()
	initCodemapAPI()
	initFileAPI()
}

// ListenAndServe start http server
func ListenAndServe(port string) error {
	return newFHServer().ListenAndServe(":" + port)
}

// ListenAndServeTLS start https server
func ListenAndServeTLS(port, certPath, keyPath string) error {
	return newFHServer().ListenAndServeTLS(":"+port, certPath, keyPath)
}

// SetHeader add api response header
func SetHeader(key string, value string) {
	httprouter.SetHeader(key, value)
}

func newFHServer() *fasthttp.Server {
	return &fasthttp.Server{
		Name:               config.Get().ServerName,
		Handler:            httprouter.FasthttpHandler(),
		MaxRequestBodySize: 5 * 1024 * 1024 * 1024 * 1024, // 5 TB
	}
}

func get(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Get(path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func post(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Post(path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func put(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Put(path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func delete(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Delete(path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func patch(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Patch(path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func head(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Head(path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func options(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Options(path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}
