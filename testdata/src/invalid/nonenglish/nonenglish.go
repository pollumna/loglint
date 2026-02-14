package nonenglish

import "log/slog"

func f() {
	logger := slog.Default()
	logger.Info("запуск сервера")                    // want `message must be in English only`
	logger.Error("ошибка подключения к базе данных") // want `message must be in English only`

	slog.Info("сервер запущен")   // want `message must be in English only`
	slog.Error("接続失敗")            // want `message must be in English only`
	slog.Warn("timeout сработал") // want `message must be in English only`
	slog.Debug("cache просчитан") // want `message must be in English only`

	slog.Info("полина " + "logged in")              // want `message must be in English only`
	logger.With().Error("done " + "запуск сервера") // want `message must be in English only`

}
