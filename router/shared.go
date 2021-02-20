package router

import "otter-cloud-ws/api/shared"

func initSharedAPI() {
	groupName := "/shared"
	var controller shared.Controller

	// Get
	get(groupName+"/folder", true, nil, controller.GetSharedFolder)
	get(groupName+"/file/list/:id/:prefix", true, nil, controller.GetObjectList)

	// Post
	post(groupName+"/add", true, nil, controller.Add)
	post(groupName+"/remove", true, nil, controller.Remove)

	// Put

	// Delete
}
