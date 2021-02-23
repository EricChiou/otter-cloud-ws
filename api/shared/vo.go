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

// GetSharedFileListReqVo request vo
type GetSharedFileListReqVo struct {
	ID     int    `json:"id" req:"true"`
	Prefix string `json:"prefix" req:"true"`
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

// GetShareableLinkReqVo request vo
type GetShareableLinkReqVo struct {
	ID             int    `json:"id" req:"true"`
	Prefix         string `json:"prefix" req:"true"`
	FileName       string `json:"fileName" req:"true"`
	ExpiresSeconds int    `json:"expiresSeconds" req:"true"`
}

// RemoveObjectReqVo request vo
type RemoveObjectReqVo struct {
	ID       int    `json:"id" req:"true"`
	Prefix   string `json:"prefix" req:"true"`
	FileName string `json:"fileName" req:"true"`
}

// UploadObjectReqVo request vo
type UploadObjectReqVo struct {
	ID     int    `json:"id" req:"true"`
	Prefix string `json:"prefix" req:"true"`
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

// GetShareableLinkResVo get shareable link request vo
type GetShareableLinkResVo struct {
	ShareableLink string `json:"shareableLink"`
}
