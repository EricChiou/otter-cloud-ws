package router

import (
	"otter-cloud-ws/acl"
	"otter-cloud-ws/api/user"
)

func initUserAPI() {
	groupName := "/user"
	var controller user.Controller

	// user list
	get(groupName+"/list", true, nil, controller.List)
	get(groupName+"/fuzzy/list", true, nil, controller.GetUserFuzzyList)

	// Post
	post(groupName+"/signUp", false, nil, controller.SignUp)
	post(groupName+"/signIn", false, nil, controller.SignIn)
	post(groupName, true, nil, controller.Update)

	// Put
	put(groupName+"/update/:userAcc", true, []acl.Code{acl.UpdateUser}, controller.UpdateByUserAcc)
	put(groupName+"/activate", true, []acl.Code{acl.UpdateUser}, controller.ActivateAcc)
}
