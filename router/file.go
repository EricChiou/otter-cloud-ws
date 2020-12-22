package router

import (
	"otter-cloud-ws/api/file"
)

func initFileAPI() {
	groupName := "/file"
	var controller file.Controller

	// Post
	post(groupName+"/list", true, nil, controller.List)
}
