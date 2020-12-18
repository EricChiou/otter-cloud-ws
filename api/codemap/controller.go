package codemap

import (
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/db/mysql"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/service/apihandler"
	"otter-cloud-ws/service/paramhandler"
)

// Controller codemap controller
type Controller struct {
	dao Dao
}

// Add add new code map
func (con *Controller) Add(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check body format
	var addReqVo AddReqVo
	if err := paramhandler.Set(webInput.Context, &addReqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	err := con.dao.Add(addReqVo)
	if err != nil {
		return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
	}

	return responseEntity.OK(ctx, nil)
}

// Update update codemap
func (con *Controller) Update(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check body format
	var updateReqVo UpdateReqVo
	if err := paramhandler.Set(webInput.Context, &updateReqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	err := con.dao.Update(updateReqVo)
	if err != nil {
		return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
	}

	return responseEntity.OK(ctx, nil)
}

// Delete delete codemap
func (con *Controller) Delete(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check param
	var deleteReqVo DeleteReqVo
	if err := paramhandler.Set(webInput.Context, &deleteReqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	err := con.dao.Delete(deleteReqVo)
	if err != nil {
		return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
	}

	return responseEntity.OK(ctx, nil)
}

// List get codemap list
func (con *Controller) List(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check param
	var listReqVo ListReqVo
	if err := paramhandler.Set(webInput.Context, &listReqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	if listReqVo.Page == 0 {
		listReqVo.Page = 1
	}
	if listReqVo.Limit == 0 {
		listReqVo.Limit = 10
	}

	list, err := con.dao.List(listReqVo)
	if err != nil {
		return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
	}

	return responseEntity.Page(ctx, list, api.Success, nil)
}
