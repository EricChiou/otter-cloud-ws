package share

import (
	"errors"
	"net/url"
	"otter-cloud-ws/api/user"
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/service/apihandler"
	"otter-cloud-ws/service/paramhandler"
	"strconv"
)

// Controller shared controller
type Controller struct {
	dao     Dao
	userDao user.Dao
}

// Add new shared folder
func (con *Controller) Add(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo AddReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	bucketName, err := con.userDao.GetBucketName(webInput.Payload.Acc)
	ownerAcc := webInput.Payload.Acc
	if err != nil {
		return responseEntity.Error(ctx, api.DBError, err)
	}

	// check folder existing
	if err := con.dao.CheckFolder(bucketName, reqVo.Prefix); err != nil {
		return responseEntity.Error(ctx, api.PrefixError, err)
	}

	// check share not duplicated
	if err := con.dao.CheckShare(ownerAcc, reqVo.SharedAcc, reqVo.Prefix); err == nil {
		return responseEntity.Error(ctx, api.Duplicate, err)
	}

	// add shared folder
	err = con.dao.Add(ownerAcc, reqVo.SharedAcc, bucketName, reqVo.Prefix, reqVo.Permission)
	if err != nil {
		return responseEntity.Error(ctx, api.DBError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// Remove shared folder
func (con *Controller) Remove(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo RemoveReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	err := con.dao.Remove(reqVo.ID, webInput.Payload.Acc)
	if err != nil {
		return responseEntity.Error(ctx, api.DBError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// Get shared folder
func (con *Controller) Get(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	getResVos := con.dao.Get(webInput.Payload.Acc)

	return responseEntity.OK(ctx, getResVos)
}

// GetObjectList by shared id
func (con *Controller) GetObjectList(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	sharedID, err := strconv.ParseInt(webInput.Context.PathParam("id"), 10, 64)
	if err != nil {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need shared id"))
	}
	prefix, err := url.QueryUnescape(webInput.Context.PathParam("prefix"))
	if err != nil {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need prefix"))
	}

	// check permission
	sharedEntity, err := con.dao.CheckPermission(int(sharedID), webInput.Payload.Acc, prefix)
	if err != nil {
		return responseEntity.Error(ctx, api.PermissionDenied, nil)
	}

	objectList := minio.ListObjects(sharedEntity.BucketName, prefix, false)

	return responseEntity.OK(ctx, objectList)
}

// Download shared file
func (con *Controller) Download(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo GetOSharedFileReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	sharedID := reqVo.ID
	prefix, _ := url.QueryUnescape(reqVo.Prefix)
	fileName, _ := url.QueryUnescape(reqVo.FileName)

	// check permission
	sharedEntity, err := con.dao.CheckPermission(int(sharedID), webInput.Payload.Acc, prefix)
	if err != nil {
		return responseEntity.Error(ctx, api.PermissionDenied, nil)
	}

	bucketName := sharedEntity.BucketName
	object, err := minio.GetObject(bucketName, prefix, fileName)
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
