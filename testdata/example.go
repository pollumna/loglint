package testdata

import (
	"fmt"
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	password := "password"
	token := "token"
	key := "key"
	username := "username"

	slog.Info("Starting server")       // Uppercase
	slog.Error("ошибка")               // Cyrillic
	slog.Info("done!🚀")                // Emoji
	slog.Info("password: " + password) // Sensitive
	slog.Info("password: 123")         // Sensitive
	slog.Info("server started")        // OK

	logger.Info("Starting server")       // Uppercase
	logger.Error("ошибка")               // Cyrillic
	logger.Info("done!🚀")                // Emoji
	logger.Info("password: 123")         // Sensitive
	logger.Info("password: " + password) // Sensitive
	logger.Info("server started")        // OK

	slog.Info("password: " + password)       // "password: "
	slog.Info(password + "password: ")       // "password: "
	slog.Info("token: " + token + " end")    // "token:  end"
	slog.Info(username + " API key: " + key) // " API key: "

	logger.Info("msg")
	logger.With().Info("msg")
	logger.With().With().Error("password")

	logger.Info(fmt.Sprintf("user password %s", password))

}
