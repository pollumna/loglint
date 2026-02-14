package sensitive

import (
	"fmt"
	"log/slog"
)

func f() {

	logger := slog.Default()

	password := "password"
	token := "token"
	key := "key"
	username := "username"

	logger.Info("user password " + password) // want `no sensitive data keywords allowed`
	logger.Debug("api_key " + key)           // want `no sensitive data keywords allowed`
	logger.Info("token " + token)            // want `no sensitive data keywords allowed`

	slog.Info("user password is 123") // want `no sensitive data keywords allowed`
	slog.Error("api_key missing")     // want `no sensitive data keywords allowed`
	slog.Warn("token expired")        // want `no sensitive data keywords allowed`
	slog.Debug("a Secret info here")  // want `no sensitive data keywords allowed`

	logger.Info(fmt.Sprintf("password reset for %s", "user")) // want `no sensitive data keywords allowed`
	logger.Info("user token: abc123")                         // want `no sensitive data keywords allowed`

	slog.With().Info("token " + token)                // want `no sensitive data keywords allowed`
	logger.With().With().Error("key" + key)           // want `no sensitive data keywords allowed`
	slog.With().With().Info("a Password " + password) // want `no sensitive data keywords allowed`

	logger.Info(fmt.Sprintf("username %s password %s", username, password)) // want `no sensitive data keywords allowed`

	slog.Info("user " + "password")        // want `no sensitive data keywords allowed`
	logger.With().Error("here " + "token") // want `no sensitive data keywords allowed`
}
