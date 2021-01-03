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

// respone vo

// GetPreviewURLResVo get preview url response vo
type GetPreviewURLResVo struct {
	URL string `json:"url"`
}
