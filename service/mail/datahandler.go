package mail

// GetMailBody to activate account
func GetMailBody(name, activeCode string) string {
	activeURL := "https://www.calicomoomoo.ml/otter-cloud/activate/" + activeCode

	html := ""
	// html += "<!DOCTYPE html>"
	// html += "<html lang=\"zh\">"

	// html += "<head>"
	// html += "  <meta charset=\"utf-8\" />"
	// html += "  <meta name=\"viewport\" content=\"width=device-width,initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no,viewport-fit=cover\">"
	// html += "  <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge,chrome=1\">"
	// html += "</head>"

	html += "<body>"
	html += "  <div>"
	html += "    <p>Hello " + name + "! Welcome to Otter Cloud.</p>"
	html += "    <p>Please click on the link to activate your account:</p>"
	html += "    <a href=\"" + activeURL + "\" target=\"_blank\">" + activeURL + "</a>"
	html += "  </div>"
	html += "</body>"

	// html += "</html>"

	return html
}
