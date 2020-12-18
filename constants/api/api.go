package api

// Authorization Bearer Token
const (
	TokenHeader string = "Authorization"
	TokenPrefix string = "Bearer "
)

// RespStatus api response status
type RespStatus string

const (
	Success          RespStatus = "ok"
	AccInactive      RespStatus = "accInactive"
	Duplicate        RespStatus = "duplicate"
	PermissionDenied RespStatus = "permissionDenied"
	TokenError       RespStatus = "tokenError"
	FormatError      RespStatus = "formatError"
	ParseError       RespStatus = "parseError"
	DBError          RespStatus = "dbError"
	DataError        RespStatus = "dataError"
	ServerError      RespStatus = "serverError"
	MinioError       RespStatus = "minioError"
	UnknownError     RespStatus = "unknownError"
)
