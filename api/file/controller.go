package file

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"otter-cloud-ws/api/shared"
	"otter-cloud-ws/api/user"
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/service/apihandler"
	"otter-cloud-ws/service/paramhandler"
	"strings"
	"time"
)

// Controller file controller
type Controller struct {
	dao       Dao
	userDao   user.Dao
	sharedDao shared.Dao
}

// List get file list
func (con *Controller) List(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var listReqVo ListReqVo
	if err := paramhandler.Set(webInput.Context, &listReqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)
	prefix, _ := url.QueryUnescape(listReqVo.Prefix)
	objectList := minio.ListObjects(bucketName, prefix, false)

	return responseEntity.OK(ctx, objectList)
}

// Upload file
func (con *Controller) Upload(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)
	prefix, _ := url.QueryUnescape(string(ctx.FormValue("prefix")))
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	prefix = strings.ReplaceAll(prefix, "//", "/")

	fileHeader, _ := ctx.FormFile("file")
	if fileHeader == nil {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need formData which key is 'file'"))
	}

	if err := minio.PutObject(bucketName, prefix, fileHeader); err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// GetPreview get object preview
func (con *Controller) GetPreview(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo GetPreviewURLReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)

	URL, err := minio.PresignedGetObject(bucketName, reqVo.Prefix, reqVo.FileName, time.Second*60*60)
	if err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	resp, err := http.Get(URL.String())
	if err != nil {
		responseEntity.Error(ctx, api.ServerError, err)
	}

	ctx.Response.Header.Add("Content-Type", "application/octet-stream")
	ctx.SetBodyStream(resp.Body, int(resp.ContentLength))

	return responseEntity.Empty()
}

// Download file
func (con *Controller) Download(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo DownloadFileReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)
	object, err := minio.GetObject(bucketName, reqVo.Prefix, reqVo.FileName)
	if err != nil {
		responseEntity.Error(ctx, api.MinioError, err)
	}

	objectInfo, err := object.Stat()
	if err != nil {
		responseEntity.Error(ctx, api.MinioError, err)
	}

	ctx.Response.Header.Add("Content-Type", "application/octet-stream")
	ctx.SetBodyStream(object, int(objectInfo.Size))

	return responseEntity.Empty()
}

// Remove file
func (con *Controller) Remove(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo RemoveFileReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)
	prefix, _ := url.QueryUnescape(reqVo.Prefix)
	fileName, _ := url.QueryUnescape(reqVo.FileName)
	err := minio.RemoveObject(bucketName, prefix, fileName)
	if err != nil {
		responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// RemoveFolder remove folder
func (con *Controller) RemoveFolder(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo RemoveFolderReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)
	prefix, _ := url.QueryUnescape(reqVo.Prefix)

	err := con.sharedDao.RemoveByPrefix(prefix, webInput.Payload.Acc)
	if err != nil {
		responseEntity.Error(ctx, api.DBError, err)
	}

	err = minio.RemoveObjects(bucketName, prefix)
	if err != nil {
		responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// GetShareableLink get object shareable link
func (con *Controller) GetShareableLink(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo GetShareableLinkReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)
	expiresSeconds := time.Duration(reqVo.ExpiresSeconds) * time.Second

	URL, err := minio.PresignedGetObject(bucketName, reqVo.Prefix, reqVo.FileName, expiresSeconds)
	if err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	resVo := GetShareableLinkResVo{
		ShareableLink: base64.StdEncoding.EncodeToString([]byte(URL.String())),
	}

	return responseEntity.OK(ctx, resVo)
}

// GetObjectByShareableLink get object shareable link
func (con *Controller) GetObjectByShareableLink(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo GetObjectByShareableLinkReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	resp, err := http.Get(reqVo.URL)
	if err != nil {
		responseEntity.Error(ctx, api.ServerError, err)
	}

	ctx.Response.Header.Add("Content-Type", "application/octet-stream")
	ctx.SetBodyStream(resp.Body, int(resp.ContentLength))

	return responseEntity.Empty()
}

// Rename file
func (con *Controller) Rename(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo RenameFileReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)
	prefix := reqVo.Prefix
	filename := reqVo.FileName
	newFilename := reqVo.NewFileName

	if _, err := minio.StatObject(bucketName, prefix+newFilename); err == nil {
		return responseEntity.Error(ctx, api.Duplicate, err)
	}

	if err := minio.RenameObject(bucketName, prefix, filename, newFilename); err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// Move files
func (con *Controller) Move(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo MoveFilesReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, _ := con.userDao.GetBucketName(webInput.Payload.Acc)

	targetFiles := minio.ListObjects(bucketName, reqVo.TargetPrefix, false)
	for _, targetFile := range targetFiles {
		if targetFile.Size > 0 && len(targetFile.ContentType) > 0 {
			targetFileNameSep := strings.SplitAfter(targetFile.Name, "/")
			targetFileName := targetFileNameSep[len(targetFileNameSep)-1]

			for _, sourceFileName := range reqVo.FileNames {
				if sourceFileName == targetFileName {
					return responseEntity.Error(ctx, api.Duplicate, nil)
				}
			}
		}
	}

	if err := minio.MoveObject(bucketName, reqVo.Prefix, reqVo.TargetPrefix, reqVo.FileNames); err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}
