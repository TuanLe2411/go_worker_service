package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"worker-service/pkg/message_system/mail"
	"worker-service/pkg/objects"

	"github.com/rs/zerolog/log"
)

func HandleUserAction(cmd any) error {
	var userActionCmd objects.UserActionCmd
	err := json.Unmarshal([]byte(cmd.(string)), &userActionCmd)
	if err != nil {
		log.Error().Str("error", fmt.Sprintf("Unmarshal user action command fail: %s, cmd: %v\n", err.Error(), cmd)).Msg("")
		return nil
	}
	link := fmt.Sprintf(os.Getenv("APP_USER_VERIFY_URL"), userActionCmd.RequestID)
	body := fmt.Sprintf("Subject: Xác nhận tài khoản\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n<html><body><p>Nhấp <a href=\"%s\">Vào Đây</a> để xác nhận.</p></body></html>", link)
	err = mail.SendEmail(userActionCmd.Email, body)
	if err != nil {
		log.Error().Str("error", fmt.Sprintf("Send email fail: %s\n", err.Error()))
		return errors.New("send email fail")
	}
	return nil
}
