package save

import (
	"log/slog"
	"net/http"
	"urlShotener/internal/lib/api/response"
	"urlShotener/internal/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Запросы поступают в виде json, который парсится в структуру Request
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// Ответ сервера
type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
	}
}
