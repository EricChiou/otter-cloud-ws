package mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"otter-cloud-ws/config/mailconfig"
)

// Send mail
func Send(mailData SendMailData) error {
	jsonBytes, err := json.Marshal(mailData)
	if err != nil {
		return err
	}

	resp, err := http.Post(mailconfig.APIURL, mailconfig.ContentType, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return errors.New("11111 " + err.Error())
	}

	var resVo Response
	json.NewDecoder(resp.Body).Decode(&resVo)
	defer resp.Body.Close()

	if resVo.Status != "ok" {
		return errors.New(resVo.Status + " " + resVo.Trace)
	}

	return nil
}
