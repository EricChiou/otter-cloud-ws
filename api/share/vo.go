package share

// request vo

// AddReqVo request vo
type AddReqVo struct {
	SharedAcc  string `json:"sharedAcc" req:"true"`
	Prefix     string `json:"prefix" req:"true"`
	Permission string `json:"permission" req:"true"`
}

// RemoveReqVo request vo
type RemoveReqVo struct {
	ID int `json:"id" req:"true"`
}

// GetOSharedFileReqVo request vo
type GetOSharedFileReqVo struct {
	ID       int    `json:"id" req:"true"`
	Prefix   string `json:"prefix" req:"true"`
	FileName string `json:"fileName" req:"true"`
}

// response vo

// GetResVo response vo
type GetResVo struct {
	ID         int    `json:"id"`
	OwnerAcc   string `json:"ownerAcc"`
	OwnerName  string `json:"ownerName"`
	Prefix     string `json:"prefix"`
	Permission string `json:"permission"`
}

// GetShareFolderResVo response vo
type GetShareFolderResVo struct {
	ID         int    `json:"id"`
	SharedAcc  string `json:"sharedAcc"`
	SharedName string `json:"sharedName"`
	Prefix     string `json:"prefix"`
	Permission string `json:"permission"`
}
