package router

import (
	"otter-cloud-ws/api/file"
)

func initFileAPI() {
	groupName := "/file"
	var controller file.Controller

	// Get
	get(groupName+"/list", true, nil, controller.List)
	get(groupName+"/preview", false, nil, controller.GetPreviewFile)

	// Post
	post(groupName+"/upload", true, nil, controller.Upload)
	post(groupName+"/download", true, nil, controller.Download)
	post(groupName+"/preview", true, nil, controller.GetPreview)
	post(groupName+"/preview/url", true, nil, controller.GetPreviewURL)
	post(groupName+"/shareableLink", true, nil, controller.GetShareableLink)
	post(groupName+"/shareableLink/object", false, nil, controller.GetObjectByShareableLink)

	// Put
	put(groupName+"/rename", true, nil, controller.Rename)
	put(groupName+"/move", true, nil, controller.Move)

	// Delete
	delete(groupName+"/remove", true, nil, controller.Remove)
	delete(groupName+"/remove/folder", true, nil, controller.RemoveFolder)
}
