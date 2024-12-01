package utils

import (
	"GoFiber-API/internal/config"
	"encoding/base64"
	"time"

	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/google/uuid"
	gopaseto "github.com/o1egl/paseto"
)

func GenerateLocalPaseto(payload string) (string, error) {
	encryptedToken, err := pasetoware.CreateToken(
		[]byte(config.GetConfig.PASETO_LOCAL_SECRET_SYMMETRIC_KEY),
		payload,
		time.Duration(config.GetConfig.PASETO_LOCAL_EXPIRATION_HOURS)*time.Hour,
		pasetoware.PurposeLocal,
	)

	if err != nil {
		return "nil", err
	}

	return encryptedToken, nil
}

func GenerateRefreshToken(userId string) string {
	return base64.StdEncoding.EncodeToString([]byte(userId + uuid.New().String()))
}

func GenerateResetPasswordToken() string {
	return base64.StdEncoding.EncodeToString([]byte(uuid.New().String()))
}

func DecryptPaseto(token string) (gopaseto.JSONToken, error) {
	var newJsonToken gopaseto.JSONToken
	var newFooter string
	err := gopaseto.NewV2().Decrypt(token, []byte(config.GetConfig.PASETO_LOCAL_SECRET_SYMMETRIC_KEY), &newJsonToken, &newFooter)

	if err != nil {
		return gopaseto.JSONToken{}, err
	}

	return newJsonToken, nil
}
