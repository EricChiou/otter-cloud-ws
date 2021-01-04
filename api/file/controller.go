package file

import (
	"encoding/base64"
	"errors"
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
		return responseEntity.Error(ctx, api.FormatError, nil)
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

	return responseEntity.OK(ctx, nil)
}

// GetPreviewURL get object preview url
func (con *Controller) GetPreviewURL(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo GetPreviewURLReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, nil)
	}

	bucketName := webInput.Payload.BucketName
	prefix, _ := url.QueryUnescape(reqVo.Prefix)
	fileName, _ := url.QueryUnescape(reqVo.FileName)

	url, err := minio.PresignedGetObject(bucketName, prefix, fileName, time.Second*60*5)
	if err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	resVo := GetPreviewURLResVo{URL: base64.StdEncoding.EncodeToString([]byte(url.String()))}
	return responseEntity.OK(ctx, resVo)
}

// Download file
func (con *Controller) Download(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo DownloadFileReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, nil)
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
		return responseEntity.Error(ctx, api.FormatError, nil)
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
		return responseEntity.Error(ctx, api.FormatError, nil)
	}

	bucketName := webInput.Payload.BucketName
	prefix, _ := url.QueryUnescape(reqVo.Prefix)
	err := minio.RemoveObjects(bucketName, prefix)
	if err != nil {
		responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}
