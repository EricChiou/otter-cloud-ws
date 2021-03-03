package jwt

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"otter-cloud-ws/config"
	"otter-cloud-ws/pkg/jwt"

	"github.com/valyala/fasthttp"
)

// Payload jwt payload struct
type Payload struct {
	Acc      string `json:"acc"`
	Name     string `json:"name"`
	RoleCode string `json:"roleCode"`
	RoleName string `json:"roleName"`
	IP       string `json:"ip,omitempty"`
	Exp      int64  `json:"exp,omitempty"`
}

// Generate generate jwt
func Generate(ctx *fasthttp.RequestCtx, acc, name, roleCode string, rememberMe bool) (string, error) {
	cfg := config.Get()
	jwtExpire, err := strconv.Atoi(cfg.JWTExpire)
	if err != nil {
		jwtExpire = 1
	}

	payload := Payload{
		Acc:      acc,
		Name:     name,
		RoleCode: roleCode,
	}
	if rememberMe {
		payload.IP, _ = GetClientIP(ctx)
	} else {
		payload.Exp = time.Now().Unix() + int64(jwtExpire*86400)
	}

	return jwt.GenerateJWT(payload, string(config.JwtAlg), cfg.JWTKey)
}

// Verify verify JWT
func Verify(jwtStr string) (Payload, bool) {
	cfg := config.Get()
	var payload Payload
	bytes, result := jwt.VerifyJWT(jwtStr, string(config.JwtAlg), cfg.JWTKey)
	if !result {
		return payload, false
	}

	if json.Unmarshal(bytes, &payload) != nil {
		return payload, false
	}

	return payload, true
}

// GetClientIP by fasthttp.RequestCtx
func GetClientIP(ctx *fasthttp.RequestCtx) (string, error) {
	clientIP := string(ctx.Request.Header.Peek("X-Forwarded-For"))
	if index := strings.IndexByte(clientIP, ','); index > -1 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP, nil
	}

	clientIP = strings.TrimSpace(string(ctx.Request.Header.Peek("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP, nil
	}

	return "", errors.New("can not get client ip")
}
