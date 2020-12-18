package interceptor

import (
	"errors"
	"otter-cloud-ws/acl"
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/service/jwt"

	"github.com/valyala/fasthttp"
)

// Acl interceptor
func Acl(ctx *fasthttp.RequestCtx, payload jwt.Payload, aclCodes ...acl.Code) error {
	if len(aclCodes) <= 0 {
		return nil
	}

	if len(payload.RoleCode) <= 0 {
		return errors.New(string(api.PermissionDenied))
	}

	// check permission
	for _, aclCode := range aclCodes {
		if ok := acl.Check(aclCode, payload.RoleCode); !ok {
			return errors.New(string(api.PermissionDenied))
		}
	}

	return nil
}
