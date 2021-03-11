package interceptor

import (
	"errors"
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/service/jwt"
	"time"

	"github.com/valyala/fasthttp"
)

// Token interceptor
func Token(ctx *fasthttp.RequestCtx) (jwt.Payload, error) {
	auth := string(ctx.Request.Header.Peek(api.TokenHeader))
	if len(auth) < len(api.TokenPrefix) {
		return jwt.Payload{}, errors.New("token error")
	}

	payload, ok := jwt.Verify(auth[len(api.TokenPrefix):])
	if !ok {
		return payload, errors.New("token error")

	}

	if payload.Exp == 0 && len(payload.IP) == 0 {
		return payload, errors.New("token error")
	}

	if payload.Exp > 0 && time.Now().Unix() > payload.Exp {
		return payload, errors.New("token error")
	}

	if len(payload.IP) > 0 {
		if clientIP, err := jwt.GetClientIP(ctx); clientIP != payload.IP || err != nil {
			return payload, errors.New("token error")
		}
	}

	return payload, nil
}
