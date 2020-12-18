package codemap

// request vo

// AddReqVo add codemap vo
type AddReqVo struct {
	Type   string `json:"type" req:"true"`
	Code   string `json:"code" req:"true"`
	Name   string `json:"name" req:"true"`
	SortNo int    `json:"sortNo"`
	Enable bool   `json:"enable"`
}

// UpdateReqVo update codemap vo
type UpdateReqVo struct {
	ID     int    `json:"id" req:"true"`
	Type   string `json:"type"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	SortNo int    `json:"sortNo"`
	Enable bool   `json:"enable"`
}

// DeleteReqVo delete codemap vo
type DeleteReqVo struct {
	ID int `json:"id" req:"true"`
}

// ListReqVo codemap list vo
type ListReqVo struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Type   string `json:"type"`
	Enable string `json:"enable"`
}

// response vo

// ListResVo codemap list data vo
type ListResVo struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	SortNo int    `json:"sortNo"`
	Enable bool   `json:"enable"`
}
