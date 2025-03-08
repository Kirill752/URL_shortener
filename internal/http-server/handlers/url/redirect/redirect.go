package redirect

import (
	"log/slog"
	"net/http"
	"urlShotener/internal/lib/api/response"
	"urlShotener/internal/lib/logger/sl"
	"urlShotener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Ответ сервера
type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLRedirector
type URLRedirector interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlRedirector URLRedirector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		// добавление параметров в логи
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		// Парсинг запроса в стрктуру Request
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("alias is empty"))
			return
		}
		// получение алиаса из БД
		URL, err := urlRedirector.GetURL(alias)
		if err != nil {
			// Проверка существоания такого алиаса
			if err == storage.ErrURLNotFound {
				log.Error("url not found", sl.Err(err), "alias", alias)
				render.JSON(w, r, response.Error("url not found"))
				return
			}
			log.Error("error while geting url", sl.Err(err))
			render.JSON(w, r, response.Error("error while geting url"))
			return
		}
		log.Info("url found", slog.String("url", URL))
		// Ответ пользователю
		http.Redirect(w, r, URL, http.StatusFound)
	}
}
