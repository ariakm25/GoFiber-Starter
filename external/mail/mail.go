package mail

import (
	"GoFiber-API/internal/config"
	internal_log "GoFiber-API/internal/log"
	"strconv"

	"gopkg.in/gomail.v2"
)

var Mail *gomail.Dialer

func NewMail() {
	host := config.GetConfig.SMTP_HOST
	port, _ := strconv.Atoi(config.GetConfig.SMTP_PORT)
	username := config.GetConfig.SMTP_USERNAME
	pass := config.GetConfig.SMTP_PASSWORD

	if host == "" || port == 0 || username == "" || pass == "" {
		internal_log.Logger.Fatal("WARN: SMTP configuration is not set")
	}

	Mail = gomail.NewDialer(host, port, username, pass)
}
