package uppercase

import "log/slog"

func f() {
	logger := slog.Default()

	slog.Info("Polina " + "logged in")          // want `message must start with lowercase letter`
	logger.Info("Starting server on port 8080") // want `message must start with lowercase letter`
	slog.Error("Failed to connect to database") // want `message must start with lowercase letter`

	slog.Info("Server started")     // want `message must start with lowercase letter`
	slog.Error("Connection failed") // want `message must start with lowercase letter`
	slog.Warn("TIMEOUT occurred")   // want `message must start with lowercase letter`
	slog.Debug("Cache miss")        // want `message must start with lowercase letter`
	logger.Info("User logged in")   // want `message must start with lowercase letter`

}
