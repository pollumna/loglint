package specialchars

import "log/slog"

func f() {

	logger := slog.Default()

	logger.Info("server started!🚀")                 // want `no special characters or emojis allowed`
	logger.Error("connection failed!!!")            // want `no special characters or emojis allowed`
	logger.Warn("warning: something went wrong...") // want `no special characters or emojis allowed`

	slog.Info("server started!!!")    // want `no special characters or emojis allowed`
	slog.Error("connection failed..") // want `no special characters or emojis allowed`
	slog.Warn("timeout occurred :)")  // want `no special characters or emojis allowed`
	slog.Debug("cache miss #1")       // want `no special characters or emojis allowed`
	slog.Info("user logged in @home") // want `no special characters or emojis allowed`

	logger.Info("something reset!") // want `no special characters or emojis allowed`

	slog.Info("polina " + "logged in!")              // want `no special characters or emojis allowed`
	logger.With().Error("done! " + "server started") // want `no special characters or emojis allowed`
}
