package valid

import (
	"fmt"
	"log/slog"
)

func goodSlogLogs() {

	logger := slog.Default()

	slog.Info("")
	fmt.Println("hello")

	slog.Info("server started")
	slog.Error("failed to connect to database")
	slog.Warn("timeout occurred")
	slog.Debug("cache miss")

	logger.Info("starting server on port 8080")
	logger.With().Info("msg")
	logger.With().With().Error("msg")

	logger.Info(fmt.Sprintf("user"))
	slog.Info(fmt.Sprintf("user: %s", "alice"))

	slog.Info("" + "text")
	slog.Info("user " + "logged in")
	logger.With().Info("user " + "logged in")

	slog.Info(("server started"))
}
