package file

// request vo

// ListReqVo get file list request vo
type ListReqVo struct {
	Prefix string `json:"prefix" req:"false"`
}

// GetPreviewURLReqVo get preview url request vo
type GetPreviewURLReqVo struct {
	Prefix   string `json:"prefix" req:"false"`
	FileName string `json:"fileName" req:"true"`
}

// DownloadFileReqVo download file request vo
type DownloadFileReqVo struct {
	Prefix   string `json:"prefix" req:"false"`
	FileName string `json:"fileName" req:"true"`
}

// RemoveFileReqVo remove file request vo
type RemoveFileReqVo struct {
	Prefix   string `json:"prefix" req:"false"`
	FileName string `json:"fileName" req:"true"`
}

// RemoveFolderReqVo remove folder requet vo
type RemoveFolderReqVo struct {
	Prefix string `json:"prefix" req:"true"`
}

// GetShareableLinkReqVo get shareable link request vo
type GetShareableLinkReqVo struct {
	Prefix         string `json:"prefix" req:"false"`
	FileName       string `json:"fileName" req:"true"`
	ContentType    string `json:"contentType" req:"true"`
	ExpiresSeconds int    `json:"expiresSeconds" req:"true"`
	ClientAddr     string `json:"clientAddr" req:"true"`
}

// GetObjectByShareableLinkReqVo get object by shareable link request vo
type GetObjectByShareableLinkReqVo struct {
	URL string `json:"url" req:"true"`
}

// respone vo

// GetPreviewURLResVo get preview url response vo
type GetPreviewURLResVo struct {
	URL string `json:"url"`
}

// GetShareableLinkResVo get shareable link request vo
type GetShareableLinkResVo struct {
	ShareableLink string `json:"shareableLink"`
}
