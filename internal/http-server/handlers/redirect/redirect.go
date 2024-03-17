package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/andrei1998Front/url-shortener/internal/lib/api/response"
	sl "github.com/andrei1998Front/url-shortener/internal/lib/logger/slog"
	"github.com/andrei1998Front/url-shortener/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Error("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		url, err := urlGetter.GetURL(alias)

		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found")
			render.JSON(w, r, response.Error("url not found"))

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
