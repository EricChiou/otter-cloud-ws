package interceptor

import (
	"otter-cloud-ws/acl"
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/service/apihandler"
	"otter-cloud-ws/service/jwt"

	"github.com/EricChiou/httprouter"
)

type WebInput struct {
	Context *httprouter.Context
	Payload jwt.Payload
}

func Set(context *httprouter.Context, needToken bool, aclCodes []acl.Code, run func(WebInput) apihandler.ResponseEntity) apihandler.ResponseEntity {
	webInput := WebInput{
		Context: context,
	}

	// check token
	payload, err := Token(context.Ctx)
	webInput.Payload = payload
	if needToken && err != nil {
		var responseEntity apihandler.ResponseEntity
		return responseEntity.Error(context.Ctx, api.TokenError, nil)
	}

	// check acl
	if err = Acl(context.Ctx, payload, aclCodes...); err != nil {
		var responseEntity apihandler.ResponseEntity
		return responseEntity.Error(context.Ctx, api.PermissionDenied, nil)
	}

	return run(webInput)
}
