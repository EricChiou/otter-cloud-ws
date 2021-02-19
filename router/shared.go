package router

import "otter-cloud-ws/api/shared"

func initSharedAPI() {
	groupName := "/shared"
	var controller shared.Controller

	// Get
	get(groupName+"/folder", true, nil, controller.GetSharedFolder)

	// Post
	post(groupName+"/add", true, nil, controller.Add)

	// Put

	// Delete
}
