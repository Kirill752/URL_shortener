package save

import (
	"errors"
	"log/slog"
	"net/http"
	"urlShotener/internal/lib/api/response"
	"urlShotener/internal/lib/logger/sl"
	"urlShotener/internal/lib/random"
	"urlShotener/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// TODO: move to config
const aliasLength = 4

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

		// Валидация структуры
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		// Если alias пустой, то генерируем его из случайных символов
		// FIXME: обработать ситуацию, когда сгенерированный alias уже встречался в таблице
		alias := req.Alias
		if alias == "" {
			alias = random.CreateRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}
		log.Info("url added", slog.Int64("id", id))
		responseOk(w, r, alias)
	}
}
func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
