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

	err := minio.CreateUserBucket(signUpData.Acc)
	if err != nil {
		return responseEntity.Error(ctx, api.MinioError, err)
	}

	activeCode := code.Get(64)
	err = con.dao.SignUp(signUpData, activeCode)
	if err != nil {
		return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
	}

	mailData := mail.SendMailData{
		From:    mail.EmailData{Name: mailconfig.FromName, Email: mailconfig.FromEmail},
		To:      []mail.EmailData{{Name: signUpData.Name, Email: signUpData.Acc}},
		Subject: mailconfig.Subject,
		Body:    mail.GetMailBody(signUpData.Name, activeCode),
	}
	if err := mail.Send(mailData); err != nil {
		return responseEntity.Error(ctx, api.ServerError, err)
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
	// check account existing and pwd
	if err != nil || signInBo.Pwd != sha3.Encrypt(signInReqVo.Pwd) {
		return responseEntity.Error(ctx, api.DataError, errors.New("incorrect account or password"))
	}

	// check account status
	if signInBo.Status == string(userstatus.Inactive) {
		return responseEntity.Error(ctx, api.AccInactive, errors.New("account inavtive"))
	}

	var signInResVo SignInResVo
	token, _ := jwt.Generate(
		ctx,
		signInBo.Acc,
		signInBo.Name,
		signInBo.RoleCode,
		signInReqVo.RememberMe,
	)
	signInResVo = SignInResVo{
		Token: token,
	}

	return responseEntity.OK(ctx, signInResVo)
}

// Update user data, PUT: /user/update
func (con *Controller) Update(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx
	payload := webInput.Payload

	// check body format
	var updateData UpdateReqVo
	if err := paramhandler.Set(webInput.Context, &updateData); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}
	if len(updateData.Name) == 0 && len(updateData.NewPwd) == 0 {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need new name or new password"))
	}

	if len(updateData.OldPwd) > 0 && len(updateData.NewPwd) > 0 {
		encryptOldPwd := sha3.Encrypt(updateData.OldPwd)
		oldPwd, err := con.dao.GetPwd(payload.Acc)

		if err != nil {
			return responseEntity.Error(ctx, mysql.ErrMsgHandler(err), err)
		}
		if encryptOldPwd != oldPwd {
			return responseEntity.Error(ctx, api.DataError, errors.New("password error"))
		}
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
	if len(updateData.Name) == 0 && len(updateData.NewPwd) == 0 {
		return responseEntity.Error(ctx, api.FormatError, errors.New("need new name or new password"))
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

// ActivateAcc by active code
func (con *Controller) ActivateAcc(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check body format
	var reqVo ActivateReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	err := con.dao.ActivateAcc(reqVo.ActiveCode)
	if err != nil {
		return responseEntity.Error(ctx, api.DataError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// SendActivationCode by account
func (con *Controller) SendActivationCode(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check body format
	var reqVo SendActivationCodeReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	activeCode := code.Get(64)
	userName, err := con.dao.SendActivationCode(reqVo.Acc, activeCode)
	if err != nil {
		return responseEntity.Error(ctx, api.DBError, err)
	}

	mailData := mail.SendMailData{
		From:    mail.EmailData{Name: mailconfig.FromName, Email: mailconfig.FromEmail},
		To:      []mail.EmailData{{Name: userName, Email: reqVo.Acc}},
		Subject: mailconfig.Subject,
		Body:    mail.GetMailBody(userName, activeCode),
	}
	if err := mail.Send(mailData); err != nil {
		return responseEntity.Error(ctx, api.ServerError, err)
	}

	return responseEntity.OK(ctx, nil)
}

// ResetPwd by account
func (con *Controller) ResetPwd(webInput interceptor.WebInput) apihandler.ResponseEntity {
	ctx := webInput.Context.Ctx

	// check body format
	var reqVo ResetPwdReqVo
	if err := paramhandler.Set(webInput.Context, &reqVo); err != nil {
		return responseEntity.Error(ctx, api.FormatError, err)
	}

	newPwd := code.Get(8)

	userName, err := con.dao.ResetPwd(reqVo.Acc, newPwd)
	if err != nil {
		return responseEntity.Error(ctx, api.DBError, err)
	}

	mailData := mail.SendMailData{
		From:    mail.EmailData{Name: mailconfig.FromName, Email: mailconfig.FromEmail},
		To:      []mail.EmailData{{Name: userName, Email: reqVo.Acc}},
		Subject: mailconfig.ResetPwdSubject,
		Body:    mail.GetResetPwdMailBody(userName, newPwd),
	}
	if err := mail.Send(mailData); err != nil {
		return responseEntity.Error(ctx, api.ServerError, err)
	}

	return responseEntity.OK(ctx, nil)
}
