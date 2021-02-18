package router

import (
	"otter-cloud-ws/api/share"
)

func initSharedAPI() {
	groupName := "/share"
	var controller share.Controller

	// Get
	get(groupName+"/folder", true, nil, controller.GetShareFolder)

	// Post
	post(groupName+"/add", true, nil, controller.Add)

	// Put

	// Delete
}
