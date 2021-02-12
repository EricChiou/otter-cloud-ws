package file

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
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
	dao Dao
}

// List get file list
func (con *Controller) List(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var listReqVo ListReqVo
	if err := paramhandler.Set(webInput.Context, &listReqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	prefix, _ := url.QueryUnescape(listReqVo.Prefix)
	objectList := minio.ListObjects(webInput.Payload.BucketName, prefix)

	return responseEntity.OK(ctx, objectList)
}

// Upload file
func (con *Controller) Upload(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	bucketName := webInput.Payload.BucketName
	prefix, _ := url.QueryUnescape(string(ctx.FormValue("prefix")))
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	fileHeader, _ := ctx.FormFile("file")
	if fileHeader == nil {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need formData which key is 'file'"))
	}

	if err := minio.PutObject(bucketName, prefix, fileHeader); err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	ctx.Response.Header.Add("Connection", "keep-alive")

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

	bucketName := webInput.Payload.BucketName
	prefix, _ := url.QueryUnescape(reqVo.Prefix)
	fileName, _ := url.QueryUnescape(reqVo.FileName)

	URL, err := minio.PresignedGetObject(bucketName, prefix, fileName, time.Second*60*60*1)
	if err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	resp, err := http.Get("http://" + URL.Host + URL.Path + "?" + URL.RawQuery)
	if err != nil {
		responseEntity.Error(ctx, api.ServerError, err)
	}

	ctx.Response.Header.Add("Connection", "keep-alive")
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

	bucketName := webInput.Payload.BucketName
	object, err := minio.GetObject(bucketName, reqVo.Prefix, reqVo.FileName)
	if err != nil {
		responseEntity.Error(ctx, api.MinioError, err)
	}

	objectInfo, err := object.Stat()
	if err != nil {
		responseEntity.Error(ctx, api.MinioError, err)
	}

	ctx.Response.Header.Add("Connection", "keep-alive")
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

	bucketName := webInput.Payload.BucketName
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

	bucketName := webInput.Payload.BucketName
	prefix, _ := url.QueryUnescape(reqVo.Prefix)
	err := minio.RemoveObjects(bucketName, prefix)
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

	bucketName := webInput.Payload.BucketName
	prefix, _ := url.QueryUnescape(reqVo.Prefix)
	fileName, _ := url.QueryUnescape(reqVo.FileName)
	contentType, _ := url.QueryUnescape(reqVo.ContentType)
	expiresSeconds := time.Duration(reqVo.ExpiresSeconds) * time.Second

	URL, err := minio.PresignedGetObject(bucketName, prefix, fileName, expiresSeconds)
	if err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	clientAddr, _ := url.QueryUnescape(reqVo.ClientAddr)
	shareableLinkURL := clientAddr + "/share-link" +
		"?fileName=" + base64.StdEncoding.EncodeToString([]byte(fileName)) +
		"&contentType=" + base64.StdEncoding.EncodeToString([]byte(contentType)) +
		"&url=" + base64.StdEncoding.EncodeToString([]byte("http://"+URL.Host+URL.Path+"?"+URL.RawQuery))

	resVo := GetShareableLinkResVo{
		ShareableLink: base64.StdEncoding.EncodeToString([]byte(shareableLinkURL)),
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

	shareableLink, err := base64.StdEncoding.DecodeString(reqVo.URL)
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.ParseError, err)
	}

	resp, err := http.Get(string(shareableLink))
	if err != nil {
		responseEntity.Error(ctx, api.ServerError, err)
	}

	ctx.Response.Header.Add("Connection", "keep-alive")
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

	bucketName := webInput.Payload.BucketName
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

	bucketName := webInput.Payload.BucketName
	prefix := reqVo.Prefix
	targetPrefix := reqVo.TargetPrefix
	filenames := reqVo.FileNames

	if err := minio.MoveObject(bucketName, prefix, targetPrefix, filenames); err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}
