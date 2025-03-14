package main

import (
	"log/slog"
	"net/http"
	"os"
	"urlShotener/internal/config"
	"urlShotener/internal/http-server/handlers/url/del"
	"urlShotener/internal/http-server/handlers/url/redirect"
	"urlShotener/internal/http-server/handlers/url/save"
	"urlShotener/internal/http-server/middleware/logger"
	"urlShotener/internal/lib/logger/sl"
	"urlShotener/internal/storage/sqllite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Примеры запросов:
// Добавление в баззу данных
// curl -X POST -H "Content-Type: application/json" -d '{"url":"https://yandex.ru", "alias":"yandex"}' -u admin:12345  http://localhost:8082/url/save
// curl -X POST -H "Content-Type: application/json" -d '{"url":"https://vk.com", "alias":"vk"}' -u admin:12345 http://localhost:8082/url/save
// Удаление из базы данных
// curl -X POST -H "Content-Type: application/json" -d '{"alias":"vk"}' -u admin:12345 http://localhost:8082/url/delete
func main() {
	// CONFIG_PATH=../../config/local.yaml go run main.go
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting url-sgortener", slog.String("env", cfg.Env))
	log.Debug("debug massages enabled")

	storage, err := sqllite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error to init storage", sl.Err(err))
		os.Exit(1)
	}
	router := chi.NewRouter()
	// middleware - это хэндлеры, которые выполняются при запуске основного хэндлера
	// Пример: если запрос модифицирующий, то помимо него нужно запустить хэндлер проверки авторизации
	// Каждому запросу будет присвоен уникальный ID
	router.Use(middleware.RequestID)
	// Получает IP подключенного клиента
	router.Use(middleware.RealIP)
	// Логируются действия клиента
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	// Востановление после паник
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/save", save.New(log, storage))
		r.Delete("/delete", del.New(log, storage))
		// r.Post("/delete", del.New(log, storage))
	})
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server error")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
