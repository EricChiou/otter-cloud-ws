package router

import (
	"otter-cloud-ws/acl"
	"otter-cloud-ws/config"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/service/apihandler"

	"github.com/EricChiou/httprouter"
	"github.com/valyala/fasthttp"
)

var mainPath string = "/otter-cloud-ws"

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
		ReadTimeout:        60 * 60 * 24 * 365,
		WriteTimeout:       60 * 60 * 24 * 365,
		IdleTimeout:        60 * 60 * 24 * 365,
	}
}

func get(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Get(mainPath+path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func post(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Post(mainPath+path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func put(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Put(mainPath+path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func delete(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Delete(mainPath+path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func patch(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Patch(mainPath+path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func head(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Head(mainPath+path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}

func options(path string, needToken bool, aclCodes []acl.Code, run func(interceptor.WebInput) apihandler.ResponseEntity) {
	httprouter.Options(mainPath+path, func(ctx *httprouter.Context) {
		interceptor.Set(ctx, needToken, aclCodes, run)
	})
}
