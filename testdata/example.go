package testdata

import (
	"fmt"
	"log/slog"
)

func main() {

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

	slog.Info("password: " + password)       // "password: "
	slog.Info(password + "password: ")       // "password: "
	slog.Info("token: " + token + " end")    // "token:  end"
	slog.Info(username + " API key: " + key) // " API key: "

	slog.Info("msg")
	slog.With().Info("msg")
	slog.With().With().Error("password")

	slog.Info(fmt.Sprintf("user password %s", password))

}
