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

	// Post
	post(groupName+"/signIn", false, nil, controller.SignIn)
	post(groupName, true, nil, controller.Update)
	post(groupName+"/:userID", true, []acl.Code{acl.UpdateUser}, controller.UpdateByUserID)

	// Put
	put(groupName+"/signUp", false, nil, controller.SignUp)

}
