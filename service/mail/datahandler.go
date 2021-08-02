package mail

// GetMailBody for activattion account
func GetMailBody(name, activeCode string) string {
	activeURL := "https://www.calicomoomoo.com/otter-cloud/activate/" + activeCode

	html := "<body>"
	html += "  <div>"
	html += "    <p>Hello " + name + "! Welcome to Otter Cloud.</p>"
	html += "    <p>Please click on the link to activate your account:</p>"
	html += "    <a href=\"" + activeURL + "\" target=\"_blank\">" + activeURL + "</a>"
	html += "  </div>"
	html += "</body>"

	return html
}

// GetResetPwdMailBody for reset pwd account
func GetResetPwdMailBody(name, newPwd string) string {
	html := "<body>"
	html += "  <div>"
	html += "    <p>Hi " + name + "! Your new password is as below.</p>"
	html += "    <p style=\"padding: 5px 8px; background-color: #eeeeee;\">" + newPwd + "</p>"
	html += "  </div>"
	html += "</body>"

	return html
}
