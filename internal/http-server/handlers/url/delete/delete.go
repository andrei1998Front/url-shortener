package delete

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/andrei1998Front/url-shortener/internal/lib/api/response"
	sl "github.com/andrei1998Front/url-shortener/internal/lib/logger/slog"
	"github.com/andrei1998Front/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Alias string `json:"alias"`
}

type Response struct {
	response.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requet_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.Alias == "" {
			log.Error("alias is empty")

			render.JSON(w, r, response.Error("invalid alias"))

			return
		}

		err = urlDeleter.DeleteURL(req.Alias)

		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Error("there are no url with this alias", slog.String("alias", req.Alias))

			render.JSON(w, r, response.Error("there are no url with this alias"))

			return
		}

		if err != nil {
			log.Error("failed to delete alias", sl.Err(err))

			render.JSON(w, r, response.Error("failed to delete alias"))

			return
		}

		log.Info("url removed")

		render.JSON(w, r, Response{Response: response.OK()})
	}
}
