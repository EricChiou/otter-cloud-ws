package file

import (
	"errors"
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/service/apihandler"
	"otter-cloud-ws/service/paramhandler"
	"strings"
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

	objectList := minio.ListObjects(webInput.Payload.BucketName, listReqVo.Prefix)

	return responseEntity.OK(ctx, objectList)
}

// Upload files
func (con *Controller) Upload(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	bucketName := webInput.Payload.BucketName
	prefix := string(ctx.FormValue("prefix"))
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	forms, _ := ctx.MultipartForm()
	if forms == nil {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need formData."))
	}

	for _, fileHeader := range forms.File["files"] {
		if err := minio.PutObject(bucketName, prefix, fileHeader); err != nil {
			return responseEntity.Error(ctx, api.MinioError, err)
		}
	}

	return responseEntity.OK(ctx, nil)
}
