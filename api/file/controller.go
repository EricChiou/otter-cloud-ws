package file

import (
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/service/apihandler"
	"otter-cloud-ws/service/paramhandler"
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
