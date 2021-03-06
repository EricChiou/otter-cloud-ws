package router

import "otter-cloud-ws/api/shared"

func initSharedAPI() {
	groupName := "/shared"
	var controller shared.Controller

	// Get
	get(groupName, true, nil, controller.GetSharedFolder)
	get(groupName+"/file/list", true, nil, controller.GetObjectList)

	// Post
	post(groupName+"/add", true, nil, controller.Add)
	post(groupName+"/remove", true, nil, controller.Remove)
	post(groupName+"/file/preview", true, nil, controller.GetPreview)
	post(groupName+"/file/preview/url", true, nil, controller.GetPreviewURL)
	post(groupName+"/file/download", true, nil, controller.Download)
	post(groupName+"/file/shareableLink", true, nil, controller.GetShareableLink)
	post(groupName+"/file/upload", true, nil, controller.UploadObject)

	// Put
	put(groupName+"/update", true, nil, controller.Update)
	put(groupName+"/file/rename", true, nil, controller.Rename)
	put(groupName+"/file/move", true, nil, controller.Move)

	// Delete
	delete(groupName+"/file", true, nil, controller.RemoveObject)
	delete(groupName+"/folder", true, nil, controller.RemoveFolder)
}
