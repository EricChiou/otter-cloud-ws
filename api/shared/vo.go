package shared

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

// GetSharedFileReqVo request vo
type GetSharedFileReqVo struct {
	ID       int    `json:"id" req:"true"`
	Prefix   string `json:"prefix" req:"true"`
	FileName string `json:"fileName" req:"true"`
}

// GetSharedFilePreviewURLReqVo request vo
type GetSharedFilePreviewURLReqVo struct {
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

// GetSharedFolderResVo response vo
type GetSharedFolderResVo struct {
	ID         int    `json:"id"`
	OwnerAcc   string `json:"ownerAcc"`
	OwnerName  string `json:"ownerName"`
	SharedAcc  string `json:"sharedAcc"`
	SharedName string `json:"sharedName"`
	Prefix     string `json:"prefix"`
	Permission string `json:"permission"`
}
