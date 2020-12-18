package user

// request vo

// SignUpReqVo user sign up request data vo
type SignUpReqVo struct {
	Acc  string `json:"acc" req:"true"`
	Pwd  string `json:"pwd" req:"true"`
	Name string `json:"name" req:"true"`
}

// SignInReqVo user sign in request data vo
type SignInReqVo struct {
	Acc string `json:"acc" req:"true"`
	Pwd string `json:"pwd" req:"true"`
}

// UpdateReqVo update user request data vo
type UpdateReqVo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

// ListReqVo List request vo
type ListReqVo struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Active string `json:"active"`
}

// response vo

// SignInResVo user sign in response data vo
type SignInResVo struct {
	Token string `json:"token"`
}

// ListResVo user list data vo
type ListResVo struct {
	ID       int    `json:"id"`
	Acc      string `json:"acc"`
	Name     string `json:"name"`
	RoleCode string `json:"roleCode"`
	Status   string `json:"status"`
}
