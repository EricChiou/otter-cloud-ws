package router

import (
	"otter-cloud-ws/api/file"
)

func initFileAPI() {
	groupName := "/file"
	var controller file.Controller

	// Get
	get(groupName+"/list", true, nil, controller.List)
	get(groupName+"/preview", true, nil, controller.GetPreview)

	// Post
	post(groupName+"/upload", true, nil, controller.Upload)
	post(groupName+"/download", true, nil, controller.Download)

	// Delete
	delete(groupName+"/remove", true, nil, controller.Remove)
	delete(groupName+"/remove/folder", true, nil, controller.RemoveFolder)
}
