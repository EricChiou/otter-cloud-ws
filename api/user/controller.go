package user

import (
	"errors"
	"net/url"
	"otter-cloud-ws/config/mailconfig"
	"otter-cloud-ws/constants/api"
	"otter-cloud-ws/constants/userstatus"
	"otter-cloud-ws/db/mysql"
	"otter-cloud-ws/interceptor"
	"otter-cloud-ws/minio"
	"otter-cloud-ws/service/apihandler"
	"otter-cloud-ws/service/code"
	"otter-cloud-ws/service/jwt"
	"otter-cloud-ws/service/mail"
	"otter-cloud-ws/service/paramhandler"
	"otter-cloud-ws/service/sha3"
)

// Controller user controller
type Controller struct {
	dao Dao
}

// SignUp user sign up controller
func (con *Controller) SignUp(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check body format
	var signUpData SignUpReqVo
	if err := paramhandler.Set(webInput.Context, &signUpData); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	if result, err := con.dao.CheckAccExisting(signUpData.Acc); result || err == nil {
		return responseEntity.Error(ctx, api.Duplicate, err)
	}

	activeCode := code.Get(64)
	mailData := mail.SendMailData{
		From:    mail.EmailData{Name: mailconfig.FromName, Email: mailconfig.FromEmail},
		To:      []mail.EmailData{{Name: signUpData.Name, Email: signUpData.Acc}},
		Subject: mailconfig.Subject,
		Body:    mail.GetMailBody(signUpData.Name, activeCode),
	}
	if err := mail.Send(mailData); err != nil {
		return responseEntity.Error(ctx, api.ServerError, err)
	}

	err := minio.CreateUserBucket(signUpData.Acc)
	if err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	err = con.dao.SignUp(signUpData, activeCode)
	if err != nil {
		return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
	}

	return responseEntity.OK(ctx, nil)
}

// SignIn user sign in controller
func (con *Controller) SignIn(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// set param
	var signInReqVo SignInReqVo
	if err := paramhandler.Set(webInput.Context, &signInReqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	signInBo, err := con.dao.SignIn(signInReqVo)
	if err != nil {
		return responseEntity.Error(ctx, api.DataError, errors.New("incorrect account or password"))
	}

	// check pwd
	if signInBo.Pwd != sha3.Encrypt(signInReqVo.Pwd) {
		return responseEntity.Error(ctx, api.DataError, errors.New("incorrect account or password"))
	}

	// check account status
	if signInBo.Status != string(userstatus.Active) {
		return responseEntity.Error(ctx, api.AccInactive, errors.New("the account is inactive"))
	}

	var signInResVo SignInResVo
	token, _ := jwt.Generate(
		signInBo.Acc,
		signInBo.Name,
		signInBo.RoleCode,
	)
	signInResVo = SignInResVo{
		Token: token,
	}

	return responseEntity.OK(ctx, signInResVo)
}

// Update user data, POST: /user
func (con *Controller) Update(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx
	payload := webInput.Payload

	// check body format
	var updateData UpdateReqVo
	if err := paramhandler.Set(webInput.Context, &updateData); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}
	if len(updateData.Name) == 0 && len(updateData.Pwd) == 0 {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need name or pwd"))
	}

	err := con.dao.Update(updateData, payload.Acc)
	if err != nil {
		return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
	}

	return responseEntity.OK(ctx, nil)
}

// UpdateByUserAcc POST: /user/:userID
func (con *Controller) UpdateByUserAcc(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	var updateData UpdateReqVo

	// check body format
	if err := paramhandler.Set(webInput.Context, &updateData); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}
	if len(updateData.Name) == 0 && len(updateData.Pwd) == 0 {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need name or pwd"))
	}

	// check path param
	userAcc := webInput.Context.PathParam("userAcc")
	if len(userAcc) == 0 {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need user account"))
	}

	err := con.dao.Update(updateData, userAcc)
	if err != nil {
		return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
	}

	return responseEntity.OK(ctx, nil)
}

// List get user list
func (con *Controller) List(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check body format
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

// GetUserFuzzyList by key word
func (con *Controller) GetUserFuzzyList(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check body format
	var reqVo FuzzyListReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	keyword, _ := url.QueryUnescape(reqVo.Keyword)
	accountList, err := con.dao.GetUserFuzzyList(keyword)
	if err != nil {
		return responseEntity.Error(ctx, api.DBError, err)
	}

	return responseEntity.OK(ctx, accountList)
}
