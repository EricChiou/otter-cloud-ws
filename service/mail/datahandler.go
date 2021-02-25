package mail

// GetMailBody to activate account
func GetMailBody(name, activeCode string) string {
	activeURL := "https://www.calicomoomoo.ml/otter-cloud/activate/" + activeCode

	html := "<body>"
	html += "  <div>"
	html += "    <p>Hello " + name + "! Welcome to Otter Cloud.</p>"
	html += "    <p>Please click on the link to activate your account:</p>"
	html += "    <a href=\"" + activeURL + "\" target=\"_blank\">" + activeURL + "</a>"
	html += "  </div>"
	html += "</body>"

	return html
}
