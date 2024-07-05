package utils

import (
	"GoFiber-API/internal/config"
	"time"

	pasetoware "github.com/gofiber/contrib/paseto"
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
