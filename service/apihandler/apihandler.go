package apihandler

import (
	"encoding/json"
	"fmt"
	"otter-cloud-ws/api/common"
	"otter-cloud-ws/constants/api"

	"github.com/valyala/fasthttp"
)

// ResponseEntity response format
type ResponseEntity struct{}

type apiResponse struct {
	Status api.RespStatus `json:"status"`
	Data   interface{}    `json:"data,omitempty"`
	Trace  interface{}    `json:"trace,omitempty"`
}

var header map[string]string = map[string]string{
	"Content-Type": "application/json",
}

func addHeader(ctx *fasthttp.RequestCtx) {
	for k, v := range header {
		ctx.Response.Header.Add(k, v)
	}
}

// OK api success
func (re *ResponseEntity) OK(ctx *fasthttp.RequestCtx, data interface{}) ResponseEntity {
	addHeader(ctx)

	result := apiResponse{
		Status: api.Success,
		Data:   data,
		Trace:  nil,
	}

	sendResp(ctx, result)
	return *re
}

// Error api error
func (re *ResponseEntity) Error(ctx *fasthttp.RequestCtx, status api.RespStatus, err error) ResponseEntity {
	addHeader(ctx)

	result := apiResponse{
		Status: status,
		Data:   nil,
		Trace:  nil,
	}

	if err != nil {
		result.Trace = err.Error()
	}

	sendResp(ctx, result)
	return *re
}

// Page api page format
func (re *ResponseEntity) Page(ctx *fasthttp.RequestCtx, list common.PageRespVo, status api.RespStatus, err error) ResponseEntity {
	addHeader(ctx)

	result := apiResponse{
		Status: status,
		Data:   list,
		Trace:  nil,
	}

	if err != nil {
		result.Trace = err.Error()
	}

	sendResp(ctx, result)
	return *re
}

func sendResp(ctx *fasthttp.RequestCtx, result apiResponse) {
	bytes, _ := json.Marshal(result)
	fmt.Fprintf(ctx, string(bytes))
}
