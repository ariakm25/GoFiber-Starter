package auth

import (
	"GoFiber-API/app/user"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/external/mail"
	"GoFiber-API/internal/config"
	internal_log "GoFiber-API/internal/log"
	"GoFiber-API/internal/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"gopkg.in/gomail.v2"
)

const (
	TypeAuthResetPasswordJob = "auth:reset-password"
)

type ResetPasswordJobPayload struct {
	Email   string
	UserUID string
}

func NewAuthResetPasswordJob(email string, userUID string) (*asynq.Task, error) {
	payload, err := json.Marshal(ResetPasswordJobPayload{Email: email, UserUID: userUID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAuthResetPasswordJob, payload), nil
}

func HandleAuthResetPasswordJob(ctx context.Context, task *asynq.Task) error {
	var payload ResetPasswordJobPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		internal_log.Logger.Error(err.Error())
		return err
	}

	token := utils.GenerateResetPasswordToken()

	userToken := &user.UserToken{
		UserID:    payload.UserUID,
		Token:     token,
		Type:      user.UserTokenTypeResetPassword,
		ExpiredAt: time.Now().Add(time.Hour * 24),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.Connection.Create(userToken).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.GetConfig.SMTP_FROM_NAME+" <"+config.GetConfig.SMTP_FROM_EMAIL+">")
	mailer.SetHeader("To", "delivered@resend.dev")
	mailer.SetHeader("Subject", "Reset Password Request")
	mailer.SetBody("text/html", "<div>Hello, you have requested to reset your password. Use this code to reset your password.</div><br/><b>"+userToken.Token+"</b>")

	err := mail.Mail.DialAndSend(mailer)

	if err != nil {
		internal_log.Logger.Error("HandleAuthResetPasswordJob Error " + err.Error())
		return err
	}

	return nil
}
