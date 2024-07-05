package auth

import (
	internal_log "GoFiber-API/internal/log"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeAuthResetPasswordJob = "auth:reset-password"
)

type ResetPasswordJobPayload struct {
	Email string `json:"email"`
}

func NewAuthResetPasswordJob(email string) (*asynq.Task, error) {
	payload, err := json.Marshal(ResetPasswordJobPayload{Email: email})
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

	// TODO Implement send an email to the user

	time.Sleep(5 * time.Second)

	log.Printf("Sending Email to User: Email=%s", payload.Email)

	return nil
}
