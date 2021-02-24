package shared

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"otter-cloud-ws/api/user"
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/constants/sharedperms"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/service/apihandler"
	"otter-cloud-ws/service/paramhandler"
	"strconv"
	"strings"
	"time"
)

// Controller shared controller
type Controller struct {
	dao     Dao
	userDao user.Dao
}

// Add new share folder
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

// GetSharedFolder by token
func (con *Controller) GetSharedFolder(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	sharedFolderList := con.dao.GetSharedFolder(webInput.Payload.Acc)

	return responseEntity.OK(ctx, sharedFolderList)
}

// Remove share folder
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

// UploadObject by shared id, prefix, fileName and token
func (con *Controller) UploadObject(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	id, err := strconv.ParseInt(string(ctx.FormValue("id")), 10, 64)
	if err != nil || id == 0 {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	prefix, err := url.QueryUnescape(string(ctx.FormValue("prefix")))
	if err != nil || len(prefix) == 0 {
		return responseEntity.Error(ctx, api.FormatError, err)
	}
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	// check permission
	sharedEntity, err := con.dao.CheckPermission(int(id), webInput.Payload.Acc, prefix)
	if err != nil || sharedEntity.Permission != sharedperms.Write {
		return responseEntity.Error(ctx, api.PermissionDenied, err)
	}

	fileHeader, _ := ctx.FormFile("file")
	if fileHeader == nil {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need formData which key is 'file'"))
	}

	if err := minio.PutObject(sharedEntity.BucketName, sharedEntity.Prefix, fileHeader); err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// RemoveObject by shared id, prefix, fileName and token
func (con *Controller) RemoveObject(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo RemoveObjectReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	prefix, _ := url.QueryUnescape(reqVo.Prefix)
	fileName, _ := url.QueryUnescape(reqVo.FileName)

	// check permission
	sharedEntity, err := con.dao.CheckPermission(reqVo.ID, webInput.Payload.Acc, prefix)
	if err != nil || sharedEntity.Permission != sharedperms.Write {
		return responseEntity.Error(ctx, api.PermissionDenied, err)
	}

	err = minio.RemoveObject(sharedEntity.BucketName, sharedEntity.Prefix, fileName)
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

	prefix, _ := url.QueryUnescape(reqVo.Prefix)

	// check permission
	sharedEntity, err := con.dao.CheckPermission(reqVo.ID, webInput.Payload.Acc, prefix)
	if err != nil || sharedEntity.Permission != sharedperms.Write {
		return responseEntity.Error(ctx, api.PermissionDenied, err)
	}

	err = minio.RemoveObjects(sharedEntity.BucketName, sharedEntity.Prefix)
	if err != nil {
		responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// GetObjectList by shared id
func (con *Controller) GetObjectList(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var fileListReqVo GetSharedFileListReqVo
	if err := paramhandler.Set(webInput.Context, &fileListReqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	prefix, _ := url.QueryUnescape(fileListReqVo.Prefix)

	// check permission
	sharedEntity, err := con.dao.CheckPermission(fileListReqVo.ID, webInput.Payload.Acc, prefix)
	if err != nil {
		return responseEntity.Error(ctx, api.PermissionDenied, err)
	}

	objectList := minio.ListObjects(sharedEntity.BucketName, sharedEntity.Prefix, false)

	return responseEntity.OK(ctx, objectList)
}

// GetPreview get shared object preview
func (con *Controller) GetPreview(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo GetSharedFilePreviewURLReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	sharedEntity, err := con.dao.CheckPermission(reqVo.ID, webInput.Payload.Acc, reqVo.Prefix)
	if err != nil {
		return responseEntity.Error(ctx, api.PermissionDenied, err)
	}

	URL, err := minio.PresignedGetObject(sharedEntity.BucketName, sharedEntity.Prefix, reqVo.FileName, time.Second*60*60)
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

// Download shared file
func (con *Controller) Download(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo GetSharedFileReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	// check permission
	sharedEntity, err := con.dao.CheckPermission(reqVo.ID, webInput.Payload.Acc, reqVo.Prefix)
	if err != nil {
		return responseEntity.Error(ctx, api.PermissionDenied, nil)
	}

	object, err := minio.GetObject(sharedEntity.BucketName, sharedEntity.Prefix, reqVo.FileName)
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

// GetShareableLink get object shareable link
func (con *Controller) GetShareableLink(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo GetShareableLinkReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	// check permission
	sharedEntity, err := con.dao.CheckPermission(reqVo.ID, webInput.Payload.Acc, reqVo.Prefix)
	if err != nil {
		return responseEntity.Error(ctx, api.PermissionDenied, nil)
	}

	expiresSeconds := time.Duration(reqVo.ExpiresSeconds) * time.Second

	URL, err := minio.PresignedGetObject(sharedEntity.BucketName, sharedEntity.Prefix, reqVo.FileName, expiresSeconds)
	if err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	resVo := GetShareableLinkResVo{
		ShareableLink: base64.StdEncoding.EncodeToString([]byte(URL.String())),
	}

	return responseEntity.OK(ctx, resVo)
}

// Rename shared file
func (con *Controller) Rename(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo RenameSharedFileReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	// check permission
	sharedEntity, err := con.dao.CheckPermission(reqVo.ID, webInput.Payload.Acc, reqVo.Prefix)
	if err != nil || sharedEntity.Permission != sharedperms.Write {
		return responseEntity.Error(ctx, api.PermissionDenied, err)
	}

	if _, err := minio.StatObject(sharedEntity.BucketName, sharedEntity.Prefix+reqVo.NewFileName); err == nil {
		return responseEntity.Error(ctx, api.Duplicate, err)
	}

	if err := minio.RenameObject(
		sharedEntity.BucketName,
		sharedEntity.Prefix,
		reqVo.FileName,
		reqVo.NewFileName,
	); err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// Move shared files
func (con *Controller) Move(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var reqVo MoveSharedFilesReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	// check sourcePrefix permission
	sourceSharedEnt, err := con.dao.CheckPermission(reqVo.ID, webInput.Payload.Acc, reqVo.Prefix)
	if err != nil || sourceSharedEnt.Permission != sharedperms.Write {
		return responseEntity.Error(ctx, api.PermissionDenied, err)
	}

	// check targetPrefix permission
	targetSharedEnt, err := con.dao.CheckPermission(reqVo.ID, webInput.Payload.Acc, reqVo.TargetPrefix)
	if err != nil || targetSharedEnt.Permission != sharedperms.Write {
		return responseEntity.Error(ctx, api.PermissionDenied, err)
	}

	targetFiles := minio.ListObjects(sourceSharedEnt.BucketName, targetSharedEnt.Prefix, false)
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

	if err := minio.MoveObject(
		sourceSharedEnt.BucketName,
		sourceSharedEnt.Prefix,
		targetSharedEnt.Prefix,
		reqVo.FileNames,
	); err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	return responseEntity.OK(ctx, nil)
}
