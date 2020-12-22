package userbo

// SignInBo sign in bo
type SignInBo struct {
	// form user table
	ID         int    `json:"id,omitempty"`
	Acc        string `json:"acc,omitempty"`
	Pwd        string `json:"pwd,omitempty"`
	Name       string `json:"name,omitempty"`
	RoleCode   string `json:"roleCode,omitempty"`
	Status     string `json:"status,omitempty"`
	BucketName string `json:"bucket_id,omitempty"`
	// from role table
	RoleName string `json:"roleName,omitempty"`
}
