package job

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"time"
)

const (
	TaskWelcome         = "email:welcome"
	TaskSignIn          = "email:signin"
	TaskPasswordChanged = "email:password_changed"
)

type WeclomeEmailPayload struct {
	To        string `json:"to"`
	FirstName string `json:"first_name"`
}

func NewWelcomeEmailTask(to, firstName string) (*asynq.Task, error) {
	payload, err := json.Marshal(WeclomeEmailPayload{
		To:        to,
		FirstName: firstName,
	})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskWelcome, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*time.Second),
	), nil
}

type SignInEmailPayload struct {
	To         string `json:"to"`
	FirstName  string `json:"first_name"`
	SignInTime string `json:"sign_in_time"`
}

func NewSignInEmailTask(to, firstName, signInTime string) (*asynq.Task, error) {
	payload, err := json.Marshal(SignInEmailPayload{
		To:         to,
		FirstName:  firstName,
		SignInTime: signInTime,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskSignIn, payload,
		asynq.MaxRetry(3),
		asynq.Queue("low"),
		asynq.Timeout(30*time.Second),
	), nil
}

type PasswordChangedEmailPayload struct {
	To        string `json:"to"`
	FirstName string `json:"first_name"`
	ChangedAt string `json:"changed_at"`
}

func NewPasswordChangedEmailTask(to, firstName, changedAt string) (*asynq.Task, error) {
	payload, err := json.Marshal(PasswordChangedEmailPayload{
		To:        to,
		FirstName: firstName,
		ChangedAt: changedAt,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskPasswordChanged, payload,
		asynq.MaxRetry(3),
		asynq.Queue("critical"),
		asynq.Timeout(30*time.Second),
	), nil
}
