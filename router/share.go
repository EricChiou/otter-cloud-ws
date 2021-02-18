package router

import (
	"otter-cloud-ws/api/share"
)

func initSharedAPI() {
	groupName := "/share"
	var controller share.Controller

	// Get

	// Post
	post(groupName+"/add", true, nil, controller.Add)

	// Put

	// Delete
}
