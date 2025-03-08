package del

import (
	"log/slog"
	"net/http"
	"urlShotener/internal/lib/api/response"
	"urlShotener/internal/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// Запросы поступают в виде json, который парсится в структуру Request
type Request struct {
	Alias string `json:"alias"`
}

// Ответ сервера
type Response struct {
	response.Response
	Alias string `json:"alias" validate:"required,alias"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) (int64, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		// Добавление параметров в логи
		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		// Парсинг запроса в JSON в структуру
		var req Request
		// err := render.DecodeForm(r.Body, &req)
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		// Валидация структуры
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		// Удаление URL из базы данных
		cnt, err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to delete url"))
			return
		}
		log.Info("url deleted", slog.Int64("number of affected url's", cnt))
		// Ответ пользователю
		responseOk(w, r, alias)
	}
}
func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
