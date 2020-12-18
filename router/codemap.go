package router

import (
	"otter-cloud-ws/acl"
	"otter-cloud-ws/api/codemap"
)

func initCodemapAPI() {
	groupName := "/codemap"
	var controller codemap.Controller

	// Get
	get(groupName+"/list", true, nil, controller.List)

	// Post
	post(groupName, true, []acl.Code{acl.UpdateCodemap}, controller.Update)

	// Put
	put(groupName, true, []acl.Code{acl.AddCodemap}, controller.Add)

	// Delete
	delete(groupName, true, []acl.Code{acl.DeleteCodemap}, controller.Delete)
}
