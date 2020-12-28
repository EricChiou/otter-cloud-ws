package router

import (
	"otter-cloud-ws/api/file"
)

func initFileAPI() {
	groupName := "/file"
	var controller file.Controller

	// Post
	post(groupName+"/list", true, nil, controller.List)

	// Put
	put(groupName+"/upload", true, nil, controller.Upload)
}
