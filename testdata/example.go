package testdata

import slog "log/slog"
import zap "go.uber.org/zap"

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	password := "password"

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

}
