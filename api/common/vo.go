package common

// PageRespVo page response vo
type PageRespVo struct {
	Records []interface{} `json:"records"`
	Page    int           `json:"page"`
	Limit   int           `json:"limit"`
	Total   int           `json:"total"`
}
