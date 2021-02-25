package mail

// SendMailData request vo
type SendMailData struct {
	From    EmailData   `json:"from"`
	To      []EmailData `json:"to"`
	Subject string      `json:"subject"`
	Body    string      `json:"body"`
}

// EmailData include Name and Email
type EmailData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Response vo
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Trace  string      `json:"trace,omitempty"`
}
